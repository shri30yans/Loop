#!/bin/bash

# Loop Project Configuration Setup Script
# This script helps set up and validate configuration across all three services

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Loop Project Configuration Setup${NC}"
echo "This script will help you set up configuration for all Loop services."
echo

# Check if services exist
if [ ! -d "Loop_backend" ] || [ ! -d "Loop_frontend" ] || [ ! -d "Loop_agent" ]; then
  echo -e "${RED}Error: One or more service directories not found.${NC}"
  echo "Please run this script from the root directory of the Loop project."
  exit 1
fi

# Function to create .env files from examples
create_env_file() {
  local service=$1
  local source_file="${service}/.env.example"
  local target_file="${service}/.env"
  
  if [ -f "$source_file" ]; then
    if [ -f "$target_file" ]; then
      echo -e "${YELLOW}Warning: ${target_file} already exists. Skipping.${NC}"
    else
      cp "$source_file" "$target_file"
      echo -e "${GREEN}Created ${target_file} from example.${NC}"
    fi
  else
    echo -e "${RED}Error: ${source_file} not found.${NC}"
  fi
}

# Create .env files for each service
echo -e "\n${BLUE}Creating .env files from examples...${NC}"
create_env_file "Loop_backend"
create_env_file "Loop_agent"

# For frontend, create .env.local instead of .env
if [ -f "Loop_frontend/.env.example" ]; then
  if [ -f "Loop_frontend/.env.local" ]; then
    echo -e "${YELLOW}Warning: Loop_frontend/.env.local already exists. Skipping.${NC}"
  else
    cp "Loop_frontend/.env.example" "Loop_frontend/.env.local"
    echo -e "${GREEN}Created Loop_frontend/.env.local from example.${NC}"
  fi
else
  echo -e "${RED}Error: Loop_frontend/.env.example not found.${NC}"
fi

# Function to prompt for configuration value
prompt_config() {
  local prompt=$1
  local default=$2
  local var_name=$3
  local is_secret=$4
  
  if [ "$is_secret" = true ]; then
    read -p "$prompt [$default]: " -s value
    echo
  else
    read -p "$prompt [$default]: " value
  fi
  
  value=${value:-$default}
  eval "$var_name='$value'"
}

echo -e "\n${BLUE}Setting up shared configuration values...${NC}"

# Prompt for shared configuration values
prompt_config "Backend URL" "http://localhost:8080" "backend_url" false
prompt_config "JWT Secret" "your_jwt_secret_here" "jwt_secret" true
prompt_config "API Key for AI Provider" "your_api_key_here" "api_key" true
prompt_config "AI Provider URL" "https://api.openai.com/v1" "api_url" false
prompt_config "AI Model Name" "gpt-4" "model_name" false

# Update backend configuration
if [ -f "Loop_backend/.env" ]; then
  echo -e "\n${BLUE}Updating Loop_backend configuration...${NC}"
  sed -i.bak "s|^JWT_SECRET=.*|JWT_SECRET=$jwt_secret|g" "Loop_backend/.env"
  sed -i.bak "s|^API_KEY=.*|API_KEY=$api_key|g" "Loop_backend/.env"
  sed -i.bak "s|^API_URL=.*|API_URL=$api_url|g" "Loop_backend/.env"
  sed -i.bak "s|^MODEL_NAME=.*|MODEL_NAME=$model_name|g" "Loop_backend/.env"
  rm "Loop_backend/.env.bak"
  echo -e "${GREEN}Updated Loop_backend/.env with shared values.${NC}"
fi

# Update agent configuration
if [ -f "Loop_agent/.env" ]; then
  echo -e "\n${BLUE}Updating Loop_agent configuration...${NC}"
  sed -i.bak "s|^GO_BACKEND_URL=.*|GO_BACKEND_URL=$backend_url|g" "Loop_agent/.env"
  sed -i.bak "s|^GO_BACKEND_API_TOKEN=.*|GO_BACKEND_API_TOKEN=$jwt_secret|g" "Loop_agent/.env"
  sed -i.bak "s|^API_KEY=.*|API_KEY=$api_key|g" "Loop_agent/.env"
  sed -i.bak "s|^API_URL=.*|API_URL=$api_url|g" "Loop_agent/.env"
  sed -i.bak "s|^MODEL_NAME=.*|MODEL_NAME=$model_name|g" "Loop_agent/.env"
  rm "Loop_agent/.env.bak"
  echo -e "${GREEN}Updated Loop_agent/.env with shared values.${NC}"
fi

# Update frontend configuration
if [ -f "Loop_frontend/.env.local" ]; then
  echo -e "\n${BLUE}Updating Loop_frontend configuration...${NC}"
  sed -i.bak "s|^NEXT_PUBLIC_BASE_API_URL=.*|NEXT_PUBLIC_BASE_API_URL=$backend_url|g" "Loop_frontend/.env.local"
  rm "Loop_frontend/.env.local.bak"
  echo -e "${GREEN}Updated Loop_frontend/.env.local with shared values.${NC}"
fi

echo -e "\n${GREEN}Configuration setup complete!${NC}"
echo "Please review the .env files in each service directory and update any additional service-specific values."
echo -e "For more information, see ${BLUE}config_documentation.md${NC}"