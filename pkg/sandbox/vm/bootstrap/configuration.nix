{ config, pkgs, lib, ... }:

{
  boot = {
    consoleLogLevel = 0;
    kernelParams = [ "quiet" "udev.log_level=3" ];

    initrd = {
      enable = true;
      verbose = false;
    };

    loader = {
      grub.enable = false;
      generationsDir = {
        enable = true;
        copyKernels = true;
      };
    };
  };

  environment = {
    # defaultPackages = [ ];
    systemPackages = with pkgs; [ curl vim ];
  };

  fileSystems = {
    "/" = {
      device = "/dev/vda";
      fsType = "ext4";
    };
  };

  nix = {
    settings = {
      auto-optimise-store = true;
      experimental-features = [ "ca-derivations" "flakes" "nix-command" "repl-flake" ];
    };
  };

  nixpkgs = {
    config = {
      allowInsecure = true;
      allowUnfree = true;
    };
    hostPlatform = lib.mkDefault "{{.System}}";
  };

  security.sudo = {
    extraConfig = "Defaults lecture = never";
    wheelNeedsPassword = false;
  };

  services.getty = {
    autologinUser = "{{.User.Username}}";
    greetingLine = "";
  };

  system.stateVersion = "23.05";

  time.timeZone = "{{.TimeZone}}";

  users.users = {
    root = {
      hashedPassword = "";
    };
    "{{.User.Username}}" = {
      isNormalUser = true;
      description = "{{.User.Name}}";
      hashedPassword = ""; # passwordless login
      extraGroups = [ "wheel" ];
    };
  };
}
