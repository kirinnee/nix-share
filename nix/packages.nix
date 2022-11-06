{ nixpkgs ? import <nixpkgs> { } }:
let
  pkgs = {
    atomi = (
      with import (fetchTarball "https://github.com/kirinnee/test-nix-repo/archive/refs/tags/v11.1.0.tar.gz");
      {
        inherit pls goland;
      }
    );

    "Unstable 6th November 2022" = (
      with import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/cae3751e9f74eea29c573d6c2f14523f41c2821a.tar.gz") { };
      {
        inherit

          coreutils
          gnugrep
          bash
          jq

          git

          go

          gocyclo
          golangci-lint
          pre-commit
          nixpkgs-fmt
          shfmt
          shellcheck;

        prettier = nodePackages.prettier;
      }
    );

  };
in
pkgs.atomi // pkgs."Unstable 6th November 2022"
