{
  description = "pd-notify";
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = inputs: with inputs; {
    overlays.default = final: prev: {
      pd-notify = prev.callPackage ./. { };
    };
  } // flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; overlays = [ self.overlays.default ]; };
    in
    {
      devShells.default = pkgs.mkShell {
        inherit (pkgs.pd-notify) CGO_ENABLED;
        buildInputs = [ pkgs.go ];
      };
      packages.default = pkgs.pd-notify;
      apps.default = { type = "app"; program = "${pkgs.pd-notify}/bin/pd-notify"; };
    });
}
