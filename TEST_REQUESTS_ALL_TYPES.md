# 🧪 TEST REQUESTS - Tất cả các Element Types

## 📋 MỤC LỤC
1. [Title Type](#1-title-type)
2. [Button Type](#2-button-type)
3. [Button URL Type](#3-button-url-type)
4. [Picture Type](#4-picture-type)
5. [Banner Type](#5-banner-type)
6. [Graphic Type](#6-graphic-type)
7. [Large Picture Type](#7-large-picture-type)
8. [Document Type](#8-document-type)
9. [Text Types](#9-text-types)
10. [Video Type](#10-video-type)
11. [Full Request - Tất cả Types](#11-full-request---tất-cả-types)

---

## 1. Title Type

**Request:**
```json
{
  "language": 1,
  "title": "Test Wiki - All Element Types",
  "keywords": "test, all types, complete",
  "public": 1,
  "level": 1,
  "unit": "Testing Unit",
  "elements": [
    {
      "number": 1,
      "type": "title",
      "value": "{\"text\":\"Welcome to Our Wiki\",\"image_key\":\"test_wiki/title_icon_123.png\",\"style\":\"large\"}"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 1,
  "type": "title",
  "value": "{\"text\":\"Welcome to Our Wiki\",...}",
  "value_json": "{\"text\":\"Welcome to Our Wiki\",...}",
  "title": {
    "text": "Welcome to Our Wiki",
    "image_key": "test_wiki/title_icon_123.png",
    "image_url": "https://cdn.example.com/test_wiki/title_icon_123.png",
    "style": "large"
  }
}
```

---

## 2. Button Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 2,
      "type": "button",
      "value": "{\"title\":\"Start Quiz\",\"code\":\"quiz_001\",\"button_icon\":\"test_wiki/btn_icon_456.png\"}"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 2,
  "type": "button",
  "value": "{\"title\":\"Start Quiz\",...}",
  "value_json": "{\"title\":\"Start Quiz\",...}",
  "button": {
    "title": "Start Quiz",
    "code": "quiz_001",
    "button_icon": "test_wiki/btn_icon_456.png",
    "button_icon_url": "https://cdn.example.com/test_wiki/btn_icon_456.png"
  }
}
```

---

## 3. Button URL Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 3,
      "type": "button_url",
      "value": "{\"title\":\"Learn More\",\"button_url\":\"https://example.com/learn\",\"button_icon\":\"test_wiki/btn_url_icon_789.png\"}"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 3,
  "type": "button_url",
  "value": "{\"title\":\"Learn More\",...}",
  "value_json": "{\"title\":\"Learn More\",...}",
  "button_url": {
    "title": "Learn More",
    "button_url": "https://example.com/learn",
    "button_icon": "test_wiki/btn_url_icon_789.png",
    "button_icon_url": "https://cdn.example.com/test_wiki/btn_url_icon_789.png"
  }
}
```

---

## 4. Picture Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 4,
      "type": "picture",
      "picture_keys": [
        {
          "key": "test_wiki/pic1_001.jpg",
          "order": 1,
          "title": "First Image"
        },
        {
          "key": "test_wiki/pic2_002.jpg",
          "order": 2,
          "title": "Second Image"
        },
        {
          "key": "test_wiki/pic3_003.jpg",
          "order": 3,
          "title": "Third Image"
        }
      ]
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 4,
  "type": "picture",
  "picture_keys": [
    {
      "key": "test_wiki/pic1_001.jpg",
      "order": 1,
      "title": "First Image"
    },
    {
      "key": "test_wiki/pic2_002.jpg",
      "order": 2,
      "title": "Second Image"
    },
    {
      "key": "test_wiki/pic3_003.jpg",
      "order": 3,
      "title": "Third Image"
    }
  ],
  "picture_keys_url": [
    {
      "order": 1,
      "url": "https://cdn.example.com/test_wiki/pic1_001.jpg"
    },
    {
      "order": 2,
      "url": "https://cdn.example.com/test_wiki/pic2_002.jpg"
    },
    {
      "order": 3,
      "url": "https://cdn.example.com/test_wiki/pic3_003.jpg"
    }
  ],
  "value_json": "[{\"key_url\":\"test_wiki/pic1_001.jpg\",\"image_url\":\"https://cdn...\",\"title\":\"First Image\",\"order\":1},{...}]"
}
```

---

## 5. Banner Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 5,
      "type": "banner",
      "value": "test_wiki/banner_hero_123.jpg"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 5,
  "type": "banner",
  "value": "test_wiki/banner_hero_123.jpg",
  "image_url": "https://cdn.example.com/test_wiki/banner_hero_123.jpg",
  "value_json": "{\"key_url\":\"test_wiki/banner_hero_123.jpg\",\"image_url\":\"https://cdn...\"}"
}
```

---

## 6. Graphic Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 6,
      "type": "graphic",
      "value": "test_wiki/graphic_diagram_456.png"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 6,
  "type": "graphic",
  "value": "test_wiki/graphic_diagram_456.png",
  "image_url": "https://cdn.example.com/test_wiki/graphic_diagram_456.png",
  "value_json": "{\"key_url\":\"test_wiki/graphic_diagram_456.png\",\"image_url\":\"https://cdn...\"}"
}
```

---

## 7. Large Picture Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 7,
      "type": "large_picture",
      "value": "test_wiki/large_pic_789.jpg"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 7,
  "type": "large_picture",
  "value": "test_wiki/large_pic_789.jpg",
  "image_url": "https://cdn.example.com/test_wiki/large_pic_789.jpg",
  "value_json": "{\"key_url\":\"test_wiki/large_pic_789.jpg\",\"image_url\":\"https://cdn...\"}"
}
```

---

## 8. Document Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 8,
      "type": "document",
      "value": "test_wiki/document_research_001.pdf"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 8,
  "type": "document",
  "value": "test_wiki/document_research_001.pdf",
  "pdf_url": "https://cdn.example.com/test_wiki/document_research_001.pdf"
}
```

---

## 9. Text Types

### 9.1 Introduction

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 9,
      "type": "introduction",
      "value": "This is an introduction to the topic. It provides a brief overview of what will be covered..."
    }
  ]
}
```

### 9.2 Main Body

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 10,
      "type": "main_body",
      "value": "The main content goes here. This section contains detailed information about the subject matter..."
    }
  ]
}
```

### 9.3 Definition

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 11,
      "type": "definition",
      "value": "Definition: A comprehensive explanation of the key terms and concepts..."
    }
  ]
}
```

### 9.4 Question

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 12,
      "type": "question",
      "value": "What are the main principles of this concept?"
    }
  ]
}
```

---

## 10. Video Type

**Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 13,
      "type": "video",
      "video_id": "video_123456789"
    }
  ]
}
```

**Response mong đợi:**
```json
{
  "number": 13,
  "type": "video",
  "video_id": "video_123456789",
  "video_url": "https://media.example.com/video_123456789.mp4"
}
```

---

## 11. Full Request - Tất cả Types

**Complete Request với tất cả element types:**

```json
{
  "language": 1,
  "title": "Complete Test Wiki - All Element Types",
  "keywords": "test, complete, all types, comprehensive",
  "public": 1,
  "image_wiki": "test_wiki/wiki_cover_main.jpg",
  "level": 1,
  "unit": "Complete Testing Unit",
  "elements": [
    {
      "number": 1,
      "type": "title",
      "value": "{\"text\":\"Complete Wiki Guide\",\"image_key\":\"test_wiki/title_icon_123.png\",\"style\":\"large\"}"
    },
    {
      "number": 2,
      "type": "banner",
      "value": "test_wiki/banner_hero_456.jpg"
    },
    {
      "number": 3,
      "type": "introduction",
      "value": "Welcome to this comprehensive guide. This wiki covers all element types available in the system..."
    },
    {
      "number": 4,
      "type": "picture",
      "picture_keys": [
        {
          "key": "test_wiki/gallery_pic1.jpg",
          "order": 1,
          "title": "First Gallery Image"
        },
        {
          "key": "test_wiki/gallery_pic2.jpg",
          "order": 2,
          "title": "Second Gallery Image"
        },
        {
          "key": "test_wiki/gallery_pic3.jpg",
          "order": 3,
          "title": "Third Gallery Image"
        }
      ]
    },
    {
      "number": 5,
      "type": "large_picture",
      "value": "test_wiki/large_feature_image.jpg"
    },
    {
      "number": 6,
      "type": "graphic",
      "value": "test_wiki/diagram_flowchart.png"
    },
    {
      "number": 7,
      "type": "main_body",
      "value": "The main content section provides detailed information. This is where the core material is presented with comprehensive explanations..."
    },
    {
      "number": 8,
      "type": "definition",
      "value": "Definition: Key terms and concepts are defined here to ensure clear understanding..."
    },
    {
      "number": 9,
      "type": "question",
      "value": "What are the key takeaways from this section?"
    },
    {
      "number": 10,
      "type": "video",
      "video_id": "video_tutorial_12345"
    },
    {
      "number": 11,
      "type": "document",
      "value": "test_wiki/research_paper.pdf"
    },
    {
      "number": 12,
      "type": "button",
      "value": "{\"title\":\"Take Quiz\",\"code\":\"quiz_final_001\",\"button_icon\":\"test_wiki/quiz_icon.png\"}"
    },
    {
      "number": 13,
      "type": "button_url",
      "value": "{\"title\":\"External Resource\",\"button_url\":\"https://example.com/resource\",\"button_icon\":\"test_wiki/external_icon.png\"}"
    }
  ]
}
```

---

## 🧪 TEST SCENARIOS

### **Scenario 1: Test Position Change**

**Initial Request:**
```json
{
  "language": 1,
  "elements": [
    { "number": 1, "type": "banner", "value": "test_wiki/banner1.jpg" },
    { "number": 2, "type": "title", "value": "{\"text\":\"Title 1\"}" },
    { "number": 3, "type": "button", "value": "{\"title\":\"Button 1\"}" }
  ]
}
```

**Update Request (Swap positions):**
```json
{
  "language": 1,
  "elements": [
    { "number": 3, "type": "button", "value": "{\"title\":\"Button 1\"}" },
    { "number": 1, "type": "banner", "value": "test_wiki/banner1.jpg" },
    { "number": 2, "type": "title", "value": "{\"text\":\"Title 1\"}" }
  ]
}
```

**Expected:** Elements giữ nguyên, chỉ đổi thứ tự, không có files bị xóa.

---

### **Scenario 2: Test Picture Order Change**

**Initial Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 1,
      "type": "picture",
      "picture_keys": [
        { "key": "pic1.jpg", "order": 1, "title": "First" },
        { "key": "pic2.jpg", "order": 2, "title": "Second" },
        { "key": "pic3.jpg", "order": 3, "title": "Third" }
      ]
    }
  ]
}
```

**Update Request (Reorder pictures):**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 1,
      "type": "picture",
      "picture_keys": [
        { "key": "pic3.jpg", "order": 1, "title": "Third" },
        { "key": "pic1.jpg", "order": 2, "title": "First" },
        { "key": "pic2.jpg", "order": 3, "title": "Second" }
      ]
    }
  ]
}
```

**Expected:** Pictures giữ nguyên, chỉ đổi order, không có files bị xóa.

---

### **Scenario 3: Test Remove Some Pictures**

**Initial Request:**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 1,
      "type": "picture",
      "picture_keys": [
        { "key": "pic1.jpg", "order": 1, "title": "First" },
        { "key": "pic2.jpg", "order": 2, "title": "Second" },
        { "key": "pic3.jpg", "order": 3, "title": "Third" }
      ]
    }
  ]
}
```

**Update Request (Keep only pic1 and pic3):**
```json
{
  "language": 1,
  "elements": [
    {
      "number": 1,
      "type": "picture",
      "picture_keys": [
        { "key": "pic1.jpg", "order": 1, "title": "First" },
        { "key": "pic3.jpg", "order": 2, "title": "Third" }
      ]
    }
  ]
}
```

**Expected:** `pic2.jpg` bị xóa, `pic1.jpg` và `pic3.jpg` giữ lại.

---

### **Scenario 4: Test Mixed Types Update**

**Initial Request:**
```json
{
  "language": 1,
  "elements": [
    { "number": 1, "type": "banner", "value": "old_banner.jpg" },
    { "number": 2, "type": "title", "value": "{\"text\":\"Old Title\"}" },
    { "number": 3, "type": "button", "value": "{\"title\":\"Old Button\"}" }
  ]
}
```

**Update Request:**
```json
{
  "language": 1,
  "elements": [
    { "number": 1, "type": "banner", "value": "new_banner.jpg" },
    { "number": 2, "type": "title", "value": "{\"text\":\"New Title\"}" },
    { "number": 4, "type": "button_url", "value": "{\"title\":\"New Button URL\"}" }
  ]
}
```

**Expected:**
- `old_banner.jpg` bị xóa, `new_banner.jpg` được thêm
- Title updated
- Element 3 (button) bị xóa
- Element 4 (button_url) được tạo mới

---

## 📊 VERIFICATION CHECKLIST

Khi test, verify những điều sau:

### **Response Structure:**
- [ ] `value` giữ nguyên key từ DB
- [ ] `value_json` chứa JSON object với URLs
- [ ] `image_url` / `pdf_url` có URLs chính xác
- [ ] `picture_keys` được sắp xếp theo `order`
- [ ] `picture_keys_url` có cùng thứ tự với `picture_keys`

### **Data Integrity:**
- [ ] Không có elements bị mất khi đổi position
- [ ] Pictures không bị xóa khi chỉ đổi order
- [ ] Files cũ bị xóa khi thay đổi thực sự
- [ ] Files được giữ khi chỉ đổi vị trí

### **File Management:**
- [ ] Files không còn dùng bị xóa khỏi AWS
- [ ] Files đang dùng không bị xóa nhầm
- [ ] Không có orphaned files

### **Edge Cases:**
- [ ] Empty picture_keys array
- [ ] Null values
- [ ] Invalid JSON trong value
- [ ] Missing files trên AWS

---

**Happy Testing! 🚀**
