#!/bin/bash

# Wiki Service Test Data Setup Script
# This script creates comprehensive test data for the Wiki Service APIs

set -e

# Configuration
BASE_URL=${BASE_URL:-"http://localhost:8000"}
AUTH_TOKEN=${AUTH_TOKEN:-"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njc4MTkzMDAsIm9yZ2FuaXphdGlvbnMiOiJTRU5CT1giLCJyb2xlcyI6IkFkbWluLCBTdHVkZW50LCBUZWFjaGVyIiwidXNlcl9pZCI6ImQ4YTc2NGIwLTJjOGUtMTFmMC1iZjhiLTFhZTBkNjVhYjViZiIsInVzZXJuYW1lIjoibXJsb2MifQ.eXVk8_76nfVHf7ReP2D7pP_mWe4zheexkgA6LypKHhQ"}

if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: AUTH_TOKEN environment variable is required"
    echo "Usage: AUTH_TOKEN=your_jwt_token_here ./setup_test_data.sh"
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

make_request() {
    local method=$1
    local url=$2
    local data=$3

    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        curl -s -X $method \
             -H "Authorization: Bearer $AUTH_TOKEN" \
             -H "Content-Type: application/json" \
             -d "$data" \
             "$BASE_URL$url"
    else
        curl -s -X $method \
             -H "Authorization: Bearer $AUTH_TOKEN" \
             "$BASE_URL$url"
    fi
}

# Step 1: Create Product Template
log_info "Step 1: Creating Product Template..."
PRODUCT_TEMPLATE_DATA='{
  "type": "product",
  "elements": [
    {
      "number": 1,
      "type": "text",
      "value": "Product Name"
    },
    {
      "number": 2,
      "type": "image",
      "value": "product_main.jpg"
    },
    {
      "number": 3,
      "type": "picture",
      "picture_keys": ["gallery1.jpg", "gallery2.jpg", "gallery3.jpg"]
    },
    {
      "number": 4,
      "type": "banner",
      "value": "banner.jpg"
    },
    {
      "number": 5,
      "type": "large_picture",
      "value": "large_banner.jpg"
    },
    {
      "number": 6,
      "type": "linked_in",
      "value": "linkedin_banner.jpg"
    },
    {
      "number": 7,
      "type": "graphic",
      "value": "graphic.jpg"
    },
    {
      "number": 8,
      "type": "document",
      "value": "manual.pdf"
    },
    {
      "number": 9,
      "type": "button",
      "value": "{\"title\":\"Buy Now\",\"button_icon\":\"cart_icon.png\"}"
    },
    {
      "number": 10,
      "type": "button_url",
      "value": "{\"title\":\"Learn More\",\"url\":\"https://example.com\",\"button_icon\":\"info_icon.png\"}"
    }
  ]
}'

response=$(make_request POST "/api/v1/wikis/template" "$PRODUCT_TEMPLATE_DATA")
if echo "$response" | grep -q '"status_code":201' || echo "$response" | grep -q '"success":true'; then
    log_success "Product template created successfully"
else
    log_error "Failed to create product template"
    echo "$response"
    exit 1
fi

# Step 2: Verify template was created
log_info "Step 2: Verifying template creation..."
response=$(make_request GET "/api/v1/wikis/template?type=product")
if (echo "$response" | grep -q '"status_code":200' || echo "$response" | grep -q '"success":true') && echo "$response" | grep -q '"type":"product"'; then
    log_success "Template verified successfully"
else
    log_error "Template verification failed"
    echo "$response"
    exit 1
fi

# Step 3: Verify 6000 wikis were created
log_info "Step 3: Verifying 6000 wikis were created..."
response=$(make_request GET "/api/v1/wikis?page=1&limit=1&type=product")
if echo "$response" | grep -q '"status_code":200' || echo "$response" | grep -q '"success":true'; then
    total=$(echo "$response" | grep -o '"total":[0-9]*' | grep -o '[0-9]*')

    if [ "$total" -ge 6000 ]; then
        log_success "6000 wikis created successfully (total: $total)"
    else
        log_error "Only $total wikis created, expected 6000+"
        exit 1
    fi
else
    log_error "Failed to get wikis list"
    echo "$response"
    exit 1
fi

# Step 4: Get first wiki ID for updates
log_info "Step 4: Getting first wiki for updates..."
response=$(make_request GET "/api/v1/wikis/code?code=0001&type=product")
wiki_id=$(echo "$response" | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"\([^"]*\)"/\1/')

if [ -n "$wiki_id" ]; then
    log_success "Got wiki ID: $wiki_id"
else
    log_error "Failed to get wiki ID"
    echo "$response"
    exit 1
fi

# Step 5: Update with English translation
log_info "Step 5: Updating with English translation..."
ENGLISH_UPDATE_DATA="{
  \"language\": 1,
  \"title\": \"Amazing Wireless Headphones\",
  \"keywords\": \"headphones, wireless, audio, technology\",
  \"level\": 3,
  \"unit\": \"pieces\",
  \"elements\": [
    {
      \"number\": 1,
      \"type\": \"text\",
      \"value\": \"Premium wireless headphones with noise cancellation\"
    },
    {
      \"number\": 2,
      \"type\": \"image\",
      \"value\": \"headphones_main_en.jpg\"
    },
    {
      \"number\": 3,
      \"type\": \"picture\",
      \"picture_keys\": [\"headphones_gallery1_en.jpg\", \"headphones_gallery2_en.jpg\", \"headphones_gallery3_en.jpg\"]
    },
    {
      \"number\": 4,
      \"type\": \"banner\",
      \"value\": \"headphones_banner_en.jpg\"
    },
    {
      \"number\": 5,
      \"type\": \"large_picture\",
      \"value\": \"headphones_large_en.jpg\"
    },
    {
      \"number\": 6,
      \"type\": \"linked_in\",
      \"value\": \"headphones_linkedin_en.jpg\"
    },
    {
      \"number\": 7,
      \"type\": \"graphic\",
      \"value\": \"headphones_graphic_en.jpg\"
    },
    {
      \"number\": 8,
      \"type\": \"document\",
      \"value\": \"headphones_manual_en.pdf\"
    },
    {
      \"number\": 9,
      \"type\": \"button\",
      \"value\": \"{\\\"title\\\":\\\"Buy Now\\\",\\\"button_icon\\\":\\\"cart_en.png\\\"}\"
    },
    {
      \"number\": 10,
      \"type\": \"button_url\",
      \"value\": \"{\\\"title\\\":\\\"Learn More\\\",\\\"url\\\":\\\"https://example.com/learn\\\",\\\"button_icon\\\":\\\"info_en.png\\\"}\"
    }
  ]
}"

response=$(make_request PUT "/api/v1/wikis/$wiki_id" "$ENGLISH_UPDATE_DATA")
if echo "$response" | grep -q '"status_code":200' || echo "$response" | grep -q '"success":true'; then
    log_success "English translation updated successfully"
else
    log_error "Failed to update English translation"
    echo "$response"
    exit 1
fi

# Step 6: Update with Vietnamese translation
log_info "Step 6: Updating with Vietnamese translation..."
VIETNAMESE_UPDATE_DATA="{
  \"language\": 2,
  \"title\": \"Tai nghe không dây cao cấp\",
  \"keywords\": \"tai nghe, không dây, âm thanh, công nghệ\",
  \"level\": 3,
  \"unit\": \"cái\",
  \"elements\": [
    {
      \"number\": 1,
      \"type\": \"text\",
      \"value\": \"Tai nghe không dây cao cấp với công nghệ chống ồn\"
    },
    {
      \"number\": 2,
      \"type\": \"image\",
      \"value\": \"headphones_main_vn.jpg\"
    },
    {
      \"number\": 3,
      \"type\": \"picture\",
      \"picture_keys\": [\"headphones_gallery1_vn.jpg\", \"headphones_gallery2_vn.jpg\", \"headphones_gallery3_vn.jpg\"]
    },
    {
      \"number\": 4,
      \"type\": \"banner\",
      \"value\": \"headphones_banner_vn.jpg\"
    },
    {
      \"number\": 5,
      \"type\": \"large_picture\",
      \"value\": \"headphones_large_vn.jpg\"
    },
    {
      \"number\": 6,
      \"type\": \"linked_in\",
      \"value\": \"headphones_linkedin_vn.jpg\"
    },
    {
      \"number\": 7,
      \"type\": \"graphic\",
      \"value\": \"headphones_graphic_vn.jpg\"
    },
    {
      \"number\": 8,
      \"type\": \"document\",
      \"value\": \"headphones_manual_vn.pdf\"
    },
    {
      \"number\": 9,
      \"type\": \"button\",
      \"value\": \"{\\\"title\\\":\\\"Mua Ngay\\\",\\\"button_icon\\\":\\\"cart_vn.png\\\"}\"
    },
    {
      \"number\": 10,
      \"type\": \"button_url\",
      \"value\": \"{\\\"title\\\":\\\"Tìm Hiểu Thêm\\\",\\\"url\\\":\\\"https://example.com/learn\\\",\\\"button_icon\\\":\\\"info_vn.png\\\"}\"
    }
  ]
}"

response=$(make_request PUT "/api/v1/wikis/$wiki_id" "$VIETNAMESE_UPDATE_DATA")
if echo "$response" | grep -q '"status_code":200' || echo "$response" | grep -q '"success":true'; then
    log_success "Vietnamese translation updated successfully"
else
    log_error "Failed to update Vietnamese translation"
    echo "$response"
    exit 1
fi

# Step 7: Update a few more wikis with different data
log_info "Step 7: Creating additional test data..."

# Update wiki 0002 with partial data
response=$(make_request PUT "/api/v1/wikis/$wiki_id" '{
  "language": 1,
  "title": "Bluetooth Speaker",
  "elements": [
    {"number": 1, "type": "text", "value": "Portable Bluetooth speaker"},
    {"number": 3, "type": "picture", "picture_keys": ["speaker1.jpg"]}
  ]
}')

# Update wiki 0003 with minimal data
WIKI_0003_RESPONSE=$(make_request GET "/api/v1/wikis/code?code=0003&type=product")
wiki_0003_id=$(echo "$WIKI_0003_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"\([^"]*\)"/\1/')

if [ -n "$wiki_0003_id" ]; then
    make_request PUT "/api/v1/wikis/$wiki_0003_id" '{
      "language": 1,
      "title": "Smart Watch",
      "elements": [
        {"number": 1, "type": "text", "value": "Fitness smartwatch"}
      ]
    }'
    log_success "Updated wiki 0003 with minimal data"
fi

# Step 8: Test statistics
log_info "Step 8: Testing statistics endpoint..."
response=$(make_request GET "/api/v1/wikis/statistics?page=1&limit=5&type=product")
if echo "$response" | grep -q '"status_code":200' || echo "$response" | grep -q '"success":true'; then
    log_success "Statistics endpoint working"
else
    log_warning "Statistics endpoint may have issues"
    echo "$response"
fi

# Step 9: Test search functionality
log_info "Step 9: Testing search functionality..."
response=$(make_request GET "/api/v1/wikis?page=1&limit=3&type=product&search=headphones")
if echo "$response" | grep -q '"status_code":200' || echo "$response" | grep -q '"success":true'; then
    log_success "Search functionality working"
else
    log_warning "Search functionality may have issues"
fi

# Summary
log_info "=== TEST DATA SETUP COMPLETE ==="
echo ""
echo "Created test data:"
echo "- Product template with 10 different element types"
echo "- 6000+ wiki instances (codes 0001-6000)"
echo "- English and Vietnamese translations for wiki 0001"
echo "- Partial updates for additional test coverage"
echo ""
echo "Wiki ID for testing: $wiki_id"
echo "Template type: product"
echo ""
echo "You can now run the Postman collection with these variables:"
echo "- base_url: $BASE_URL"
echo "- auth_token: [SET IN POSTMAN ENVIRONMENT]"
echo "- wiki_id: $wiki_id"
echo "- template_type: product"
