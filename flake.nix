{
  description = "pd-notify";
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";
  outputs = inputs: with inputs;
    let
      forAllSystems = f: nixpkgs.lib.genAttrs
        [ "aarch64-linux" "x86_64-linux" "aarch64-darwin" "x86_64-darwin" ]
        (system: f {
          inherit system;
          pkgs = import nixpkgs { inherit system; overlays = [ self.overlays.default ]; };
        });
    in
    {
      overlays.default = final: prev: { pd-notify = prev.callPackage ./. { }; };
      packages = forAllSystems ({ pkgs, ... }: {
        default = pkgs.pd-notify;
      });
      apps = forAllSystems ({ pkgs, ... }: {
        default = { type = "app"; program = "${pkgs.pd-notify}/bin/pd-notify"; };
      });
      devShells = forAllSystems ({ pkgs, ... }: {
        default = pkgs.mkShell {
          inherit (pkgs.pd-notify) CGO_ENABLED;
          buildInputs = [ pkgs.pd-notify.go ];
        };
      });
    };
}
