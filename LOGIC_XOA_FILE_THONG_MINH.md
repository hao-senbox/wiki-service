# 🧠 LOGIC XÓA FILE THÔNG MINH - Chỉ xóa khi KEY thực sự thay đổi

## 🎯 MỤC TIÊU

Đảm bảo files chỉ bị xóa khi KEY thực sự thay đổi, KHÔNG xóa khi chỉ đổi vị trí (number).

## 📋 YÊU CẦU

1. **Đổi vị trí (position change):** Banner/Graphic/Picture ở element number 1 → number 2 với CÙNG key → **KHÔNG XÓA**
2. **Thay đổi key thực sự:** Banner thay từ `old.jpg` → `new.jpg` → **XÓA `old.jpg`**
3. **Xóa element:** Element bị xóa khỏi request và key không còn dùng → **XÓA FILE**

## 🔧 GIẢI PHÁP ĐÃ IMPLEMENT

### **Bước 1: Thu thập tất cả file keys trong REQUEST**

```go
// Collect all file keys being used in the request (to avoid deleting files that are just repositioned)
requestFileKeys := make(map[string]bool)
for _, reqElem := range reqElements {
    // Collect single file keys (banner, graphic, etc.)
    if reqElem.Value != nil && *reqElem.Value != "" {
        requestFileKeys[*reqElem.Value] = true
    }
    // Collect picture keys
    for _, picItem := range reqElem.PictureKeys {
        if picItem.Key != "" {
            requestFileKeys[picItem.Key] = true
        }
    }
}
```

**Mục đích:** Tạo một "whitelist" chứa TẤT CẢ file keys đang được sử dụng trong request hiện tại.

### **Bước 2: Kiểm tra trước khi xóa file**

```go
// Check if old file is still being used in the request
if !requestFileKeys[*existingElem.Value] {
    // File not in request anymore - safe to delete
    u.fileGateway.DeleteImage(ctx, *existingElem.Value)
}
// If file is still in request, skip deletion (just position change)
```

**Logic:** Chỉ xóa file nếu nó KHÔNG có trong `requestFileKeys` (không còn được dùng nữa).

## 📊 CÁC SCENARIO CHI TIẾT

### **Scenario 1: Đổi vị trí banner (KHÔNG XÓA)**

**DB ban đầu:**
```json
{
  "elements": [
    { "number": 1, "type": "banner", "value": "banner.jpg" },
    { "number": 2, "type": "text", "value": "Nội dung" }
  ]
}
```

**Request:**
```json
{
  "elements": [
    { "number": 5, "type": "banner", "value": "banner.jpg" },  // Đổi number 1→5
    { "number": 2, "type": "text", "value": "Nội dung" }
  ]
}
```

**Logic xử lý:**
1. `requestFileKeys = {"banner.jpg": true}`
2. Process element 5 (banner mới):
   - `existingElem` = element number 1 (banner cũ)
   - `*existingElem.Value = "banner.jpg"`
   - Check: `requestFileKeys["banner.jpg"] = true` → **KHÔNG XÓA**
3. Cleanup element 1:
   - Element 1 không có trong request
   - Check: `requestFileKeys["banner.jpg"] = true` → **KHÔNG XÓA**

**✅ Kết quả:** File `banner.jpg` được giữ lại

---

### **Scenario 2: Thay đổi banner thực sự (CÓ XÓA)**

**DB ban đầu:**
```json
{
  "elements": [
    { "number": 1, "type": "banner", "value": "cu_banner.jpg" }
  ]
}
```

**Request:**
```json
{
  "elements": [
    { "number": 1, "type": "banner", "value": "moi_banner.jpg" }  // Thay key
  ]
}
```

**Logic xử lý:**
1. `requestFileKeys = {"moi_banner.jpg": true}`
2. Process element 1:
   - `*existingElem.Value = "cu_banner.jpg"`
   - Check: `requestFileKeys["cu_banner.jpg"] = false` → **XÓA FILE**

**✅ Kết quả:** File `cu_banner.jpg` bị xóa, giữ `moi_banner.jpg`

---

### **Scenario 3: Xóa element banner (CÓ XÓA)**

**DB ban đầu:**
```json
{
  "elements": [
    { "number": 1, "type": "banner", "value": "banner.jpg" },
    { "number": 2, "type": "text", "value": "Nội dung" }
  ]
}
```

**Request:**
```json
{
  "elements": [
    { "number": 2, "type": "text", "value": "Nội dung" }  // Không có banner
  ]
}
```

**Logic xử lý:**
1. `requestFileKeys = {}` (không có file nào)
2. Cleanup element 1:
   - Element 1 không có trong request
   - Check: `requestFileKeys["banner.jpg"] = false` → **XÓA FILE**

**✅ Kết quả:** File `banner.jpg` bị xóa

---

### **Scenario 4: Pictures - Đổi thứ tự (KHÔNG XÓA)**

**DB ban đầu:**
```json
{
  "elements": [
    {
      "number": 1,
      "type": "picture",
      "picture_keys": [
        {"key": "pic1.jpg", "order": 1},
        {"key": "pic2.jpg", "order": 2},
        {"key": "pic3.jpg", "order": 3}
      ]
    }
  ]
}
```

**Request:**
```json
{
  "elements": [
    {
      "number": 1,
      "type": "picture",
      "picture_keys": [
        {"key": "pic1.jpg", "order": 1},
        {"key": "pic3.jpg", "order": 2}  // Giữ pic1, pic3, bỏ pic2
      ]
    }
  ]
}
```

**Logic xử lý:**
1. `requestFileKeys = {"pic1.jpg": true, "pic3.jpg": true}`
2. Process picture keys:
   - Check keys to delete: `["pic2.jpg"]`
   - Check: `requestFileKeys["pic2.jpg"] = false` → **XÓA `pic2.jpg`**
   - Check: `requestFileKeys["pic1.jpg"] = true` → **KHÔNG XÓA**
   - Check: `requestFileKeys["pic3.jpg"] = true` → **KHÔNG XÓA**

**✅ Kết quả:** Chỉ xóa `pic2.jpg`, giữ `pic1.jpg` và `pic3.jpg`

---

### **Scenario 5: Swap positions của 2 banners (KHÔNG XÓA)**

**DB ban đầu:**
```json
{
  "elements": [
    { "number": 1, "type": "banner", "value": "banner1.jpg" },
    { "number": 2, "type": "banner", "value": "banner2.jpg" }
  ]
}
```

**Request:**
```json
{
  "elements": [
    { "number": 2, "type": "banner", "value": "banner1.jpg" },  // Swap
    { "number": 1, "type": "banner", "value": "banner2.jpg" }   // Swap
  ]
}
```

**Logic xử lý:**
1. `requestFileKeys = {"banner1.jpg": true, "banner2.jpg": true}`
2. Process cả 2 elements:
   - Check: `requestFileKeys["banner1.jpg"] = true` → **KHÔNG XÓA**
   - Check: `requestFileKeys["banner2.jpg"] = true` → **KHÔNG XÓA**

**✅ Kết quả:** Cả 2 files đều được giữ lại

---

## 🔍 ĐIỂM MẠNH CỦA LOGIC MỚI

### **1. Simplicity (Đơn giản)**
- Chỉ cần 1 map `requestFileKeys` để track files đang dùng
- Logic clear: File trong request → giữ, không trong request → xóa

### **2. Correctness (Chính xác)**
- Xử lý đúng TẤT CẢ scenarios: position change, key change, delete element
- Không bị nhầm lẫn giữa position change và key change

### **3. Performance (Hiệu năng)**
- O(1) lookup trong map
- Chỉ scan request 1 lần để build map
- Không có nested loops phức tạp

### **4. Coverage (Bao phủ)**
- Xử lý tất cả types: banner, graphic, linked_in, large_picture, file, picture
- Xử lý cả single file và picture arrays

---

## ⚠️ EDGE CASES ĐÃ XỬ LÝ

### **Edge Case 1: File được dùng ở nhiều elements**
- File `shared.jpg` dùng ở element 1 và element 2
- Xóa element 1 nhưng element 2 vẫn dùng
- **Kết quả:** File không bị xóa (vì vẫn trong requestFileKeys)

### **Edge Case 2: Empty request**
- Request không có elements nào
- **Kết quả:** Tất cả files bị xóa (đúng behavior)

### **Edge Case 3: Duplicate file keys trong request**
- Request có 2 elements cùng dùng 1 file
- **Kết quả:** File chỉ được add vào map 1 lần, logic vẫn đúng

### **Edge Case 4: Null/Empty values**
- Element có value = null hoặc ""
- **Kết quả:** Không add vào requestFileKeys, logic vẫn đúng

---

## 📈 SO SÁNH LOGIC CŨ VS MỚI

| **Aspect** | **Logic Cũ** | **Logic Mới** |
|---|---|---|
| **Position change** | ❌ Xóa nhầm file | ✅ Giữ file |
| **Key change** | ✅ Xóa đúng | ✅ Xóa đúng |
| **Delete element** | ✅ Xóa đúng | ✅ Xóa đúng |
| **Complexity** | O(n²) nested loops | O(n) single pass |
| **Edge cases** | ❌ Nhiều bugs | ✅ Xử lý đầy đủ |
| **Maintainability** | ❌ Khó hiểu | ✅ Dễ hiểu |

---

## 🎯 KẾT LUẬN

**Logic mới đảm bảo:**
- ✅ Chỉ xóa khi KEY thực sự thay đổi
- ✅ KHÔNG xóa khi chỉ đổi vị trí
- ✅ Xử lý đúng tất cả edge cases
- ✅ Performance tối ưu O(n)
- ✅ Code dễ hiểu và maintain

**Bug "thỉnh thoảng mất file khi đổi vị trí" đã được fix hoàn toàn! 🎉**
