package main

import "fmt"

func printCompletions() {
	// Zsh completion function that reads cached topics.
	// Install: cht --completions > ~/.zsh/completions/_cht
	// Or:      eval "$(cht --completions)"
	fmt.Print(`#compdef cht

_cht() {
  local cache_file="${HOME}/.cache/cht/topics"

  if [[ ! -f "$cache_file" ]]; then
    cht --refresh 2>/dev/null
  fi

  if [[ -f "$cache_file" ]]; then
    compadd -f -X "cheat.sh topics" -- "${(@f)$(< "$cache_file")}"
  fi
}

_cht "$@"
`)
}
