#!/usr/bin/env bash
# Bootstrap a new microservice project from this template.
#
# Replaces:
#   - Go module path (go.mod + imports)
#   - ShortMicroserviceName (= project name without leading 5s-)
#   - GraphQL example field prefix tmp_* → <short_name with _>_*
#
# Usage:
#   make init_project
#   make init_project MODULE=5s-dataset-api GQL_SUFFIX=dts
#
# Prompts:
#   1) MODULE     — project name (example: 5s-dataset-api)
#   2) GQL_SUFFIX — short Federation field prefix instead of tmp_
#                   (supergraph style: dts, acc, chat, exch → dts_eventGet, acc_getCache, …)
#
# ShortMicroserviceName is derived automatically: 5s-dataset-api → dataset-api.
# If MODULE has no '/', it becomes github.com/testovoleg/<MODULE>.
#
# Optional flags: SKIP_GQLGEN=1 DRY_RUN=1 FORCE=1

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

OLD_MODULE="github.com/testovoleg/5s-microservice-template"
OLD_SHORT="microservice-template"
OLD_GQL_PREFIX="tmp"
MODULE_PREFIX="github.com/testovoleg"

MODULE="${MODULE:-}"
GQL_SUFFIX="${GQL_SUFFIX:-}"
SKIP_GQLGEN="${SKIP_GQLGEN:-0}"
DRY_RUN="${DRY_RUN:-0}"
FORCE="${FORCE:-0}"

log()  { printf '%s\n' "$*"; }
die()  { printf 'error: %s\n' "$*" >&2; exit 1; }

to_pascal() {
  local s="${1//-/_}"
  local out="" part
  IFS='_' read -ra parts <<< "$s"
  for part in "${parts[@]}"; do
    [[ -z "$part" ]] && continue
    out+="$(printf '%s' "${part:0:1}" | tr '[:lower:]' '[:upper:]')${part:1}"
  done
  printf '%s' "$out"
}

# prompt_var VAR "prompt text" [default] [required=0|1]
# Reads from /dev/tty when possible so prompts work under `make`.
prompt_var() {
  local var_name="$1" prompt="$2" default="${3:-}" required="${4:-0}"
  local value="${!1:-}"
  local input_src=""

  if [[ -n "$value" ]]; then
    printf -v "$var_name" '%s' "$value"
    log "Using $var_name=$value"
    return 0
  fi

  if [[ -r /dev/tty ]]; then
    input_src=/dev/tty
  elif [[ -t 0 ]]; then
    input_src=""
  else
    if [[ -n "$default" ]]; then
      printf -v "$var_name" '%s' "$default"
      log "Using $var_name=$default (non-interactive default)"
      return 0
    fi
    die "$var_name is required (non-interactive shell)"
  fi

  local hint="$prompt"
  if [[ -n "$default" ]]; then
    hint+=" [$default]"
  fi

  if [[ -n "$input_src" ]]; then
    read -r -p "$hint: " value <"$input_src"
  else
    read -r -p "$hint: " value
  fi

  if [[ -z "$value" ]]; then
    value="$default"
  fi
  if [[ -z "$value" && "$required" == "1" ]]; then
    die "$var_name is required"
  fi
  printf -v "$var_name" '%s' "$value"
}

current_module() {
  awk '/^module / { print $2; exit }' go.mod
}

replace_in_file() {
  local file="$1" from="$2" to="$3"
  if [[ "$from" == "$to" ]]; then
    return 0
  fi
  if ! grep -qF -- "$from" "$file" 2>/dev/null; then
    return 0
  fi
  if [[ "$DRY_RUN" == "1" ]]; then
    log "DRY: replace in $file: '$from' → '$to'"
    return 0
  fi
  FROM="$from" TO="$to" perl -i -0pe 's/\Q$ENV{FROM}\E/$ENV{TO}/g' "$file"
}

list_text_files() {
  find . -type f \
    \( \
      -name '*.go' -o -name '*.yaml' -o -name '*.yml' -o -name '*.md' -o \
      -name '*.graphqls' -o -name '*.proto' -o -name 'go.mod' -o -name 'Makefile' -o \
      -name 'Dockerfile' -o -name '*.Dockerfile' -o -name 'Dockerfile.*' \
    \) \
    ! -path './.git/*' \
    ! -path './vendor/*' \
    ! -path './node_modules/*' \
    ! -path './.idea/*' \
    ! -path './.vscode/*'
}

# --- resolve inputs ---

prompt_var MODULE "1/2 MODULE (project name, e.g. 5s-dataset-api)" "" 1

# Allow short form: 5s-dataset-api → github.com/testovoleg/5s-dataset-api
if [[ "$MODULE" != */* ]]; then
  MODULE="${MODULE_PREFIX}/${MODULE}"
fi

PROJECT_NAME="${MODULE##*/}"
SHORT_NAME="${PROJECT_NAME#5s-}"

# Default GQL suffix: short name without trailing -api, dashes → underscore.
# 5s-dataset-api → dataset (prefer a short Federation prefix like dts/acc/chat).
default_gql="${SHORT_NAME%-api}"
default_gql="${default_gql//-/_}"
log "GQL_SUFFIX — короткий префикс полей в supergraph вместо tmp_ (как dts, acc, chat → dts_eventGet)."
prompt_var GQL_SUFFIX "2/2 GQL_SUFFIX (e.g. dts)" "$default_gql" 1

# Normalize: dataset-api → dataset_api; used as field prefix (tmp_tipicalData → dataset_tipicalData)
GQL_SUFFIX="${GQL_SUFFIX//-/_}"

if [[ "$MODULE" == "$OLD_MODULE" ]]; then
  die "MODULE must differ from template module ($OLD_MODULE)"
fi

if [[ ! "$MODULE" =~ ^[a-zA-Z0-9._/-]+$ ]]; then
  die "MODULE contains invalid characters: $MODULE"
fi

if [[ ! "$SHORT_NAME" =~ ^[a-zA-Z0-9][a-zA-Z0-9_-]*$ ]]; then
  die "derived ShortMicroserviceName is invalid: $SHORT_NAME (from $PROJECT_NAME)"
fi

if [[ ! "$GQL_SUFFIX" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
  die "GQL_SUFFIX must be a valid GraphQL name: $GQL_SUFFIX"
fi

cur="$(current_module)"
if [[ "$cur" != "$OLD_MODULE" && "$FORCE" != "1" ]]; then
  die "go.mod module is already '$cur' (expected '$OLD_MODULE'). Re-run with FORCE=1 to continue."
fi

OLD_GQL_PASCAL="$(to_pascal "$OLD_GQL_PREFIX")"
NEW_GQL_PASCAL="$(to_pascal "$GQL_SUFFIX")"

log ""
log "Bootstrap project"
log "  MODULE:                 $OLD_MODULE → $MODULE"
log "  ShortMicroserviceName:  $OLD_SHORT → $SHORT_NAME"
if [[ "$GQL_SUFFIX" != "$OLD_GQL_PREFIX" ]]; then
  log "  GQL_SUFFIX:             ${OLD_GQL_PREFIX}_* → ${GQL_SUFFIX}_*  (Go: ${OLD_GQL_PASCAL}* → ${NEW_GQL_PASCAL}*)"
else
  log "  GQL_SUFFIX:             unchanged ($GQL_SUFFIX)"
fi
[[ "$DRY_RUN" == "1" ]] && log "  mode: DRY_RUN"
log ""

# --- rewrite files ---

while IFS= read -r file; do
  replace_in_file "$file" "$OLD_MODULE" "$MODULE"
done < <(list_text_files)

replace_in_file "pkg/constants/constants.go" \
  "ShortMicroserviceName = \"${OLD_SHORT}\"" \
  "ShortMicroserviceName = \"${SHORT_NAME}\""

if [[ "$GQL_SUFFIX" != "$OLD_GQL_PREFIX" ]]; then
  gql_files=(
    ./graphql_service/schema/queries.graphqls
    ./graphql_service/schema/mutations.graphqls
    ./graphql_service/internal/app/delivery/graphql/queries.resolvers.go
    ./graphql_service/internal/app/delivery/graphql/mutations.resolvers.go
    ./graphql_service/schema/generated.go
  )
  for file in "${gql_files[@]}"; do
    [[ -f "$file" ]] || continue
    replace_in_file "$file" "tmp_tipicalData" "${GQL_SUFFIX}_tipicalData"
    replace_in_file "$file" "tmp_tipicalMutation" "${GQL_SUFFIX}_tipicalMutation"
    replace_in_file "$file" "tmp_eventGet" "${GQL_SUFFIX}_eventGet"
    replace_in_file "$file" "TmpTipicalData" "${NEW_GQL_PASCAL}TipicalData"
    replace_in_file "$file" "TmpTipicalMutation" "${NEW_GQL_PASCAL}TipicalMutation"
    replace_in_file "$file" "TmpEventGet" "${NEW_GQL_PASCAL}EventGet"
    replace_in_file "$file" "Mutation.tmp_tipicalMutation" "Mutation.${GQL_SUFFIX}_tipicalMutation"
    replace_in_file "$file" "Query.tmp_tipicalData" "Query.${GQL_SUFFIX}_tipicalData"
    replace_in_file "$file" "Query.tmp_eventGet" "Query.${GQL_SUFFIX}_eventGet"
    replace_in_file "$file" "field_Mutation_tmp_tipicalMutation" "field_Mutation_${GQL_SUFFIX}_tipicalMutation"
    replace_in_file "$file" "field_Query_tmp_tipicalData" "field_Query_${GQL_SUFFIX}_tipicalData"
    replace_in_file "$file" "field_Query_tmp_eventGet" "field_Query_${GQL_SUFFIX}_eventGet"
    replace_in_file "$file" "_Mutation_tmp_tipicalMutation" "_Mutation_${GQL_SUFFIX}_tipicalMutation"
    replace_in_file "$file" "_Query_tmp_tipicalData" "_Query_${GQL_SUFFIX}_tipicalData"
    replace_in_file "$file" "_Query_tmp_eventGet" "_Query_${GQL_SUFFIX}_eventGet"
  done

  if [[ "$SKIP_GQLGEN" != "1" ]]; then
    if [[ "$DRY_RUN" == "1" ]]; then
      log "DRY: make gqlgen"
    else
      log "Regenerating GraphQL (make gqlgen)..."
      make gqlgen
    fi
  else
    log "Skipping gqlgen (SKIP_GQLGEN=1). Run: make gqlgen"
  fi
fi

if [[ "$DRY_RUN" != "1" ]]; then
  log "Running go mod tidy..."
  go mod tidy
  log "Verifying build..."
  go build ./...
fi

log ""
log "Done. Next steps:"
log "  1. Review git diff (module path, ShortMicroserviceName, GraphQL tmp_* rename)"
log "  2. Adjust */config/config.yaml (ports, kafka topics, OTL, external APIs)"
log "  3. Rename tipical/example domain stubs to real commands when ready"
log "  4. Update README for the new service"
