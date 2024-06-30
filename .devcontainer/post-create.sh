#!/usr/bin/env bash
sudo apt-get install make --yes
grep --quiet --fixed-strings --line-regexp 'source .devcontainer/git-completion.bash' ~/.bashrc || echo 'source .devcontainer/git-completion.bash' >> ~/.bashrc
curl -q -JLO https://github.com/dandavison/delta/releases/download/0.17.0/git-delta_0.17.0_amd64.deb && sudo dpkg -i git-delta_0.17.0_amd64.deb
