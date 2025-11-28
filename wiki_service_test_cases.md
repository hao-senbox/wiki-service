# Wiki Service - Comprehensive Test Cases & Postman Data

## Overview
Wiki Service là một hệ thống quản lý wiki đa ngôn ngữ với các tính năng:
- Template management (tạo 6000 wiki instances)
- Multi-language support (1: English, 2: Vietnamese)
- File management (images, PDFs, buttons)
- Statistics reporting
- Search functionality

## API Endpoints

### 1. POST /api/v1/wikis/template - Create Wiki Template
**Purpose**: Tạo template và 6000 wiki instances cho một type cụ thể

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "type": "product",
  "elements": [
    {
      "number": 1,
      "type": "text",
      "value": "Product description"
    },
    {
      "number": 2,
      "type": "image",
      "value": "product_image_key.jpg"
    },
    {
      "number": 3,
      "type": "picture",
      "picture_keys": ["pic1.jpg", "pic2.jpg", "pic3.jpg"]
    },
    {
      "number": 4,
      "type": "button",
      "value": "{\"title\":\"Buy Now\",\"button_icon\":\"cart_icon.png\"}"
    },
    {
      "number": 5,
      "type": "button_url",
      "value": "{\"title\":\"Learn More\",\"url\":\"https://example.com\",\"button_icon\":\"info_icon.png\"}"
    }
  ]
}
```

### 2. GET /api/v1/wikis/template - Get Template
**Query Parameters**:
- `type` (required): Type of template

### 3. GET /api/v1/wikis/statistics - Get Statistics
**Query Parameters**:
- `page` (default: 1)
- `limit` (default: 20)
- `type` (optional): Filter by type
- `search` (optional): Search by code or title

### 4. GET /api/v1/wikis/code - Get Wiki by Code
**Query Parameters**:
- `code` (required): Wiki code (0001-6000)
- `language` (optional): Language filter (1 or 2)
- `type` (required): Wiki type

### 5. GET /api/v1/wikis - Get Wikis List
**Query Parameters**:
- `page` (default: 1)
- `limit` (default: 20)
- `language` (optional): Language filter
- `type` (required): Wiki type
- `search` (optional): Search by code or title

### 6. GET /api/v1/wikis/:id - Get Wiki by ID
**Path Parameters**:
- `id`: MongoDB ObjectID

**Query Parameters**:
- `language` (optional): Language filter

### 7. PUT /api/v1/wikis/:id - Update Wiki
**Path Parameters**:
- `id`: MongoDB ObjectID

**Request Body**:
```json
{
  "language": 1,
  "title": "Updated English Title",
  "elements": [
    {
      "number": 1,
      "type": "text",
      "value": "Updated description"
    },
    {
      "number": 2,
      "type": "picture",
      "picture_keys": ["new_pic1.jpg", "new_pic2.jpg"]
    }
  ]
}
```

## Test Scenarios

### Scenario 1: Template Creation & Management
1. **Create Product Template**
2. **Verify Template Created** - GET /template?type=product
3. **Verify 6000 Wikis Created** - GET /wikis?type=product&page=1&limit=10
4. **Check First Wiki** - GET /wikis/code?code=0001&type=product

### Scenario 2: Multi-Language Support
1. **Update English Translation**
2. **Update Vietnamese Translation**
3. **Get Wiki with English Filter**
4. **Get Wiki with Vietnamese Filter**
5. **Get Wiki without Language Filter**

### Scenario 3: File Management
1. **Update with New Images**
2. **Update with New Picture Keys**
3. **Verify Old Files Deleted**
4. **Update Button Icons**

### Scenario 4: Search Functionality
1. **Search by Code**
2. **Search by Title**
3. **Search with Special Characters**
4. **Search Multiple Results**

### Scenario 5: Statistics
1. **Get Statistics for All Wikis**
2. **Get Statistics with Search**
3. **Verify Statistics Accuracy**
4. **Test Pagination**

### Scenario 6: Error Handling
1. **Invalid Parameters**
2. **Missing Required Fields**
3. **Invalid IDs**
4. **Unauthorized Access**

## Test Data Setup

### Step 1: Create Template
```bash
curl -X POST http://localhost:8023/api/v1/wikis/template \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "product",
    "elements": [
      {"number": 1, "type": "text", "value": "Product Name"},
      {"number": 2, "type": "image", "value": "product_main.jpg"},
      {"number": 3, "type": "picture", "picture_keys": ["gallery1.jpg", "gallery2.jpg"]},
      {"number": 4, "type": "banner", "value": "banner.jpg"},
      {"number": 5, "type": "button", "value": "{\"title\":\"Buy\",\"button_icon\":\"cart.png\"}"},
      {"number": 6, "type": "button_url", "value": "{\"title\":\"Details\",\"url\":\"http://example.com\",\"button_icon\":\"info.png\"}"},
      {"number": 7, "type": "document", "value": "manual.pdf"}
    ]
  }'
```

### Step 2: Update Translations
```bash
# Update English (language 1)
curl -X PUT http://localhost:8023/api/v1/wikis/WIKI_ID_HERE \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "language": 1,
    "title": "Amazing Product",
    "keywords": "product, amazing, buy",
    "level": 3,
    "unit": "pieces",
    "elements": [
      {"number": 1, "type": "text", "value": "Amazing Product Description"},
      {"number": 2, "type": "image", "value": "product_main_en.jpg"},
      {"number": 3, "type": "picture", "picture_keys": ["gallery1_en.jpg", "gallery2_en.jpg"]},
      {"number": 4, "type": "banner", "value": "banner_en.jpg"}
    ]
  }'

# Update Vietnamese (language 2)
curl -X PUT http://localhost:8023/api/v1/wikis/WIKI_ID_HERE \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "language": 2,
    "title": "Sản phẩm Tuyệt vời",
    "keywords": "sản phẩm, tuyệt vời, mua",
    "level": 3,
    "unit": "cái",
    "elements": [
      {"number": 1, "type": "text", "value": "Mô tả sản phẩm tuyệt vời"},
      {"number": 2, "type": "image", "value": "product_main_vn.jpg"},
      {"number": 3, "type": "picture", "picture_keys": ["gallery1_vn.jpg", "gallery2_vn.jpg"]},
      {"number": 4, "type": "banner", "value": "banner_vn.jpg"}
    ]
  }'
```

## Element Types Reference

| Type | Description | Value Format | Notes |
|------|-------------|--------------|-------|
| text | Plain text | String | |
| image | Single image | File key | Auto-generates image_url |
| banner | Banner image | File key | Auto-generates image_url |
| large_picture | Large image | File key | Auto-generates image_url |
| linked_in | LinkedIn image | File key | Auto-generates image_url |
| graphic | Graphic image | File key | Auto-generates image_url |
| document | PDF document | File key | Auto-generates pdf_url |
| picture | Multiple images | Array of file keys | Auto-generates picture_keys_url |
| button | Button with icon | JSON string | `{"title":"...", "button_icon":"..."}` |
| button_url | Button with URL | JSON string | `{"title":"...", "url":"...", "button_icon":"..."}` |

## Response Formats

### Wiki Response
```json
{
  "success": true,
  "message": "Wiki fetched successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "code": "0001",
    "public": 1,
    "translation": [
      {
        "language": 1,
        "title": "Amazing Product",
        "keywords": "product, amazing",
        "level": 3,
        "unit": "pieces",
        "elements": [
          {
            "number": 1,
            "type": "text",
            "value": "Product description"
          },
          {
            "number": 2,
            "type": "image",
            "value": "product_main.jpg",
            "image_url": "https://cdn.example.com/product_main.jpg"
          },
          {
            "number": 3,
            "type": "picture",
            "picture_keys": ["pic1.jpg", "pic2.jpg"],
            "picture_keys_url": ["https://cdn.example.com/pic1.jpg", "https://cdn.example.com/pic2.jpg"]
          },
          {
            "number": 4,
            "type": "button",
            "button": {
              "title": "Buy Now",
              "button_icon": "cart.png",
              "button_icon_url": "https://cdn.example.com/cart.png"
            }
          }
        ]
      }
    ],
    "image_wiki": "",
    "creator": {
      "id": "user123",
      "username": "john_doe",
      "nickname": "John",
      "fullname": "John Doe",
      "email": "john@example.com",
      "avatar": "https://cdn.example.com/avatar.jpg"
    },
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

### Statistics Response
```json
{
  "success": true,
  "message": "Statistics fetched successfully",
  "data": [
    {
      "languages": {
        "1": "yes",
        "2": "yes"
      },
      "code": "0001",
      "elements": [
        {
          "number": 1,
          "type": "text",
          "check": "yes"
        },
        {
          "number": 2,
          "type": "image",
          "check": "yes"
        },
        {
          "number": 3,
          "type": "picture",
          "check": "no"
        }
      ]
    }
  ]
}
```

## Error Responses

### Common Error Codes
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Missing/invalid token
- `404 Not Found`: Wiki/template not found
- `500 Internal Server Error`: Server errors

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```
