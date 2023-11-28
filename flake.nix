{
  description = "pd-notify";
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";
  outputs = inputs: with inputs; {
    overlays.default = final: prev: { pd-notify = prev.callPackage ./. { }; };
    legacyPackages = nixpkgs.lib.genAttrs
      [ "aarch64-linux" "x86_64-linux" "aarch64-darwin" "x86_64-darwin" ]
      (system: import nixpkgs { inherit system; overlays = [ self.overlays.default ]; });
    apps = nixpkgs.lib.mapAttrs
      (_: pkgs: {
        default = { type = "app"; program = "${pkgs.pd-notify}/bin/pd-notify"; };
      })
      self.legacyPackages;
    devShells = nixpkgs.lib.mapAttrs
      (_: pkgs: {
        default = pkgs.mkShell {
          inherit (pkgs.pd-notify) nativeBuildInputs;
          buildInputs = with pkgs; [ revive ];
        };
      })
      self.legacyPackages;
  };
}

