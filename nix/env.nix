{ nixpkgs ? import <nixpkgs> { } }:
let pkgs = import ./packages.nix { inherit nixpkgs; }; in
with pkgs;
{
  system = [
    coreutils
    gnugrep
    bash
    jq
  ];

  dev = [
    goland
  ];

  main = [
    pls
    git
    go
  ];

  lint = [
    gocyclo
    golangci-lint
    pre-commit
    nixpkgs-fmt
    prettier
    shfmt
    shellcheck
  ];

  ops = [
  ];

}
