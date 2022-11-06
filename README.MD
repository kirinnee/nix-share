# Nix Share

Allows automatically adding all machine with nix-store within the same subnet

This is achieved by sending UDP packets to all known machine within the subnet with `sx-go` periodically. Each machine
is also with the client to receive and process these UDP packet, updating their `subtituters` and `trusted-public-keys`,
and discarding these `subtituters` and `trusted-public-keys` if too much time has passed.

# Setting up

--TODO--
