#!/usr/bin/env bash
# agents/codex.sh — OpenAI Codex CLI runtime adapter.
#
# Uses the OpenAI Codex CLI: https://github.com/openai/codex
# Install: npm install -g @openai/codex
# Auth: codex login (ChatGPT account) or set OPENAI_API_KEY.

set -euo pipefail

codex_model() {
    case "${AGENT_MODEL:-default}" in
        default|codex) echo "gpt-5-codex" ;;
        gpt-5)         echo "gpt-5" ;;
        *)             echo "$AGENT_MODEL" ;;
    esac
}

run_agent_cycle() {
    local model
    model="$(codex_model)"
    echo "[codex] cycle starting with model: $model"

    # In the unattended launcher, stdin stays attached to the container process.
    # Redirect from /dev/null so Codex doesn't wait for extra input before starting.
    codex exec \
        --model "$model" \
        --sandbox workspace-write \
        "$AGENT_PROMPT" \
        </dev/null
}

check_agent_installed() {
    command -v codex >/dev/null 2>&1 || {
        echo "ERROR: codex CLI not found. Install: npm install -g @openai/codex"
        return 1
    }
}

check_agent_authed() {
    if [ -z "${OPENAI_API_KEY:-}" ] && [ ! -f "$HOME/.codex/auth.json" ]; then
        echo "ERROR: Codex not authed. Run: codex login   (or set OPENAI_API_KEY)"
        return 1
    fi
}
