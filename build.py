#!/usr/bin/env python3
# -*- coding: utf-8 -*-

#
# Copyright (c) 2018 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import argparse
import os
import os.path
import re
import subprocess
import sys

# The import path of the project:
IMPORT_PATH = "github.com/container-mgmt/messaging-library"

# The name and version of the project:
PROJECT_NAME = "messaging-library"
PROJECT_VERSION = "0.0.0"

# The list of tools that will be executed directly inside the Go environment
# but without parsing the command line. For example, if the script is invoked
# with these command line:
#
#   build dep ensure -v
#
# It will execute 'dep ensure -v' inside the go environment but it will not
# parse or process the 'ensure' and '-v' options.
DIRECT_TOOLS = [
  "dep",
  "go",
]

# The values extracted from the command line:
argv = None


def say(what):
    """
    Writes a message to the standard output, and then flushes it, so that
    the output doesn't appear out of order.
    """
    print(what, flush=True)


def cache(function):
    """
    A decorator that creates a cache for the results of a function, so that
    when the function is called the second time the result will be returned
    from the cache without actually executing it.
    """
    cache = dict()

    def helper(*key):
        try:
            value = cache[key]
        except KeyError:
            value = function(*key)
            cache[key] = value
        return value
    return helper


def find_paths(base, include=None, exclude=None):
    """
    Recursively finds the paths inside the 'base' directory whose complete
    names match the 'include' regular expression and don't match the 'exclude'
    regular expression. By default all paths are included and no path is
    excluded.
    """
    include_re = re.compile(include) if include else None
    exclude_re = re.compile(exclude) if exclude else None
    paths = []
    for root, _, names in os.walk(base):
        for name in names:
            path = os.path.join(root, name)
            path = os.path.abspath(path)
            should_include = include_re and include_re.search(path)
            should_exclude = exclude_re and exclude_re.search(path)
            if should_include and not should_exclude:
                paths.append(path)
    return paths


def go_tool(*args):
    """
    Executes a command with the 'GOPATH' environment variable pointing to the
    project specific path, and with the symbolic link as the working directory.
    The argument should be a list containing the complete command line to
    execute.
    """
    # Make sure that the required directories exist:
    go_path = ensure_go_path()
    project_link = ensure_project_link()

    # Modify the environment so that the Go tool will find the project files
    # using the `GOPATH` environment variable. Note that setting the `PWD`
    # environment is necessary, because the `cwd` is always resolved to a
    # real path by the operating system.
    env = dict(os.environ)
    env["GOPATH"] = go_path
    env["PWD"] = project_link

    # Run the Go tool and wait till it finishes:
    say("Running command '{args}'".format(args=" ".join(args)))
    process = subprocess.Popen(
        args=args,
        env=env,
        cwd=project_link,
    )
    result = process.wait()
    if result != 0:
        raise Exception("Command '{args}' failed with exit code {code}".format(
            args=" ".join(args),
            code=result,
        ))


@cache
def ensure_project_dir():
    say("Calculating project directory")
    return os.path.dirname(os.path.realpath(__file__))


@cache
def ensure_go_path():
    """
    Creates and returns the '.gopath' directory that will be used as the
    'GOPATH' for the project.
    """
    project_dir = ensure_project_dir()
    go_path = os.path.join(project_dir, '.gopath')
    if not os.path.exists(go_path):
        say('Creating Go path `{path}`'.format(path=go_path))
        os.mkdir(go_path)
    return go_path


@cache
def ensure_go_bin():
    """
    Creates and returns the Go 'bin' directory that will be used for the
    project.
    """
    go_path = ensure_go_path()
    go_bin = os.path.join(go_path, "bin")
    if not os.path.exists(go_bin):
        os.mkdir(go_bin)
    return go_bin


@cache
def ensure_go_pkg():
    """
    Creates and returns the Go 'pkg' directory that will be used for the
    project.
    """
    go_path = ensure_go_path()
    go_pkg = os.path.join(go_path, "pkg")
    if not os.path.exists(go_pkg):
        os.mkdir(go_pkg)
    return go_pkg


@cache
def ensure_go_src():
    """
    Creates and returns the Go 'src' directory that will be used for the
    project.
    """
    go_path = ensure_go_path()
    go_src = os.path.join(go_path, "src")
    if not os.path.exists(go_src):
        os.mkdir(go_src)
    return go_src


@cache
def ensure_src_link(import_path, src_path):
    """
    Creates the symbolik link that will be used to make the source for the
    given import path appear in the 'GOPATH' expected by go tools. Returns the
    full path of the link.
    """
    go_src = ensure_go_src()
    src_link = os.path.join(go_src, import_path)
    link_dir = os.path.dirname(src_link)
    if not os.path.exists(link_dir):
        os.makedirs(link_dir)
    if not os.path.exists(src_link):
        os.symlink(src_path, src_link)
    return src_link


@cache
def ensure_project_link():
    """
    Creates the symbolik link that will be used to make the project appear
    in the 'GOPATH' expected by go tools. Returns the full path of the link.
    """
    project_dir = ensure_project_dir()
    return ensure_src_link(IMPORT_PATH, project_dir)


@cache
def ensure_vendor_dir():
    """
    Creates and populates the 'vendor' directory if it doesn't exist yet.
    Returns the full path of the directory.
    """
    project_link = ensure_project_link()
    vendor_dir = os.path.join(project_link, "vendor")
    if not os.path.exists(vendor_dir):
        go_tool("dep", "ensure", "--vendor-only", "-v")
    return vendor_dir


@cache
def ensure_package_paths():
    """
    Returns the list of import paths of the packages of the project.
    """
    # Create an import path for each subdirectory of the 'pkg' directory that
    # contains at least one Go source file:
    project_dir = ensure_project_dir()
    pkg_dir = os.path.join(project_dir, "pkg")
    pkg_paths = []
    for root, _, names in os.walk(pkg_dir):
        if any([name.endswith(".go") for name in names]):
            rel_dir = os.path.relpath(root, project_dir)
            rel_path = rel_dir.replace(os.sep, "/")
            pkg_path = IMPORT_PATH + "/" + rel_path
            pkg_paths.append(pkg_path)

    # Sort the import paths, for predictable behaviour:
    pkg_paths.sort()

    return pkg_paths


@cache
def ensure_packages():
    """
    Builds all the packages of the project.
    """
    # Make sure that the vendor directory is populated:
    ensure_vendor_dir()

    # Build the packages:
    pkg_paths = ensure_package_paths()
    go_tool("go", "install", *pkg_paths)


@cache
def ensure_binaries():
    """
    Builds the binaries corresponding to each subdirectory of the 'cmd'
    directory. Returns a list containing the absolute path names of the
    generated binaries.
    """
    # Make sure that the vendor directory is populated:
    ensure_vendor_dir()

    # Get the names of the subdirectories of the 'cmd' directory:
    project_dir = ensure_project_dir()
    cmd_dir = os.path.join(project_dir, 'cmd')
    cmd_names = []
    for cmd_name in os.listdir(cmd_dir):
        cmd_path = os.path.join(cmd_dir, cmd_name)
        if os.path.isdir(cmd_path):
            cmd_names.append(cmd_name)
    cmd_names.sort()

    # Build the binaries:
    for cmd_name in cmd_names:
        say("Building binary '{name}'".format(name=cmd_name))
        cmd_path = "{path}/cmd/{name}".format(path=IMPORT_PATH, name=cmd_name)
        go_tool("go", "install", cmd_path)

    # Build the result:
    go_bin = ensure_go_bin()
    result = []
    for cmd_name in cmd_names:
        cmd_path = os.path.join(go_bin, cmd_name)
        result.append(cmd_path)
    return result


def build():
    """
    Implements the 'packages' subcommand.
    """
    ensure_packages()


def build_binaries():
    """
    Implements the 'binaries' subcommand.
    """
    ensure_binaries()


def test():
    """
    Runs the unit tests.
    """
    pkg_paths = ensure_package_paths()
    go_tool("go", "test", *pkg_paths)


def bench():
    """
    Runs the benchmarks.
    """
    pkg_paths = ensure_package_paths()
    go_tool("go", "test", "-bench=.", *pkg_paths)


def lint():
    """
    Runs the 'golint' tool on all the source files.
    """
    go_tool(
        "golint",
        "-min_confidence", "0.9",
        "-set_exit_status",
        "./pkg/...",
        "./cmd/...")


def fmt():
    """
    Formats all the source files of the project using the 'gofmt' tool.
    """
    go_tool("gofmt", "-s", "-l", "-w", "./pkg/", "./cmd/")


def main():
    # Create the top level command line parser:
    parser = argparse.ArgumentParser(
        prog=os.path.basename(sys.argv[0]),
        description="A simple build tool, just for this project.",
    )
    parser.add_argument(
        "--verbose",
        help="Genenerate verbose output.",
        default=False,
        action="store_true",
    )
    parser.add_argument(
        "--debug",
        help="Genenerate debug, very verbose, output.",
        default=False,
        action="store_true",
    )
    subparsers = parser.add_subparsers()

    # Create the parser for the 'build' command:
    build_parser = subparsers.add_parser("build")
    build_parser.set_defaults(func=build)

    # Create the parser for the 'test' command:
    test_parser = subparsers.add_parser("test")
    test_parser.set_defaults(func=test)

    # Create the parser for the 'bench' command:
    bench_parser = subparsers.add_parser("bench")
    bench_parser.set_defaults(func=bench)

    # Create the parser for the 'lint' command:
    lint_parser = subparsers.add_parser("lint")
    lint_parser.set_defaults(func=lint)

    # Create the parser for the 'fmt' command:
    fmt_parser = subparsers.add_parser("fmt")
    fmt_parser.set_defaults(func=fmt)

    # Create the parser for the 'binaries' command:
    binaries_parser = subparsers.add_parser("binaries")
    binaries_parser.set_defaults(func=build_binaries)

    # Run the selected tool:
    code = 0
    if len(sys.argv) > 0 and sys.argv[1] in DIRECT_TOOLS:
        go_tool(*sys.argv[1:])
    else:
        global argv
        argv = parser.parse_args()
        if not hasattr(argv, "func"):
            parser.print_usage()
            code = 1
        argv.func()

    # Bye:
    sys.exit(code)


if __name__ == "__main__":
    main()
