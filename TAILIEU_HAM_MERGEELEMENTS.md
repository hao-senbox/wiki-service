# Tài liệu kỹ thuật - Hàm mergeElements trong Wiki Service

## Mục lục
1. [Tổng quan hàm](#tổng-quan-hàm)
2. [Cấu trúc code](#cấu-trúc-code)
3. [Logic xử lý chi tiết](#logic-xử-lý-chi-tiết)
4. [Quản lý file](#quản-lý-file)
5. [Xử lý lỗi](#xử-lý-lỗi)
6. [Hiệu năng](#hiệu-năng)
7. [Các trường hợp test](#các-trường-hợp-test)
8. [Trường hợp đặc biệt](#trường-hợp-đặc-biệt)

## Tổng quan hàm

```go
func (u *wikiUseCase) mergeElements(ctx context.Context, translation *entity.Translation, reqElements []request.Element) error
```

### Mục đích
Hàm `mergeElements` có nhiệm vụ hợp nhất (merge) các elements từ HTTP request vào bản dịch wiki hiện có, xử lý việc tạo mới, cập nhật và xóa elements cùng với việc dọn dẹp file liên quan.

### Trách nhiệm chính
- **Quản lý elements**: Thêm mới, cập nhật hoặc xóa elements trong wiki
- **Dọn dẹp file**: Tự động xóa hình ảnh/PDF không còn sử dụng
- **Đảm bảo tính nhất quán**: Giữ cho database luôn chính xác sau mỗi thao tác
- **Quản lý khóa ảnh**: Xử lý mảng khóa ảnh phức tạp với thứ tự sắp xếp

### Tham số
- `ctx context.Context`: Context của request để xử lý timeout và hủy bỏ
- `translation *entity.Translation`: Bản dịch hiện tại cần chỉnh sửa
- `reqElements []request.Element`: Các elements mới từ HTTP request

### Giá trị trả về
- `error`: Trả về lỗi nếu có vấn đề khi xóa file hoặc xử lý

---

## Cấu trúc code

### 1. Giai đoạn chuẩn bị (Dòng 634-642)

```go
// Tạo map để tra cứu elements hiện có theo số thứ tự
existingElements := make(map[int]*entity.Element)
for i := range translation.Elements {
    elem := &translation.Elements[i]
    existingElements[elem.Number] = elem
}

// Tạo map để theo dõi elements nào được cập nhật
updatedNumbers := make(map[int]bool)
```

**Mục đích**: Tạo cấu trúc dữ liệu hiệu quả để tra cứu elements với độ phức tạp O(1).

**Cấu trúc dữ liệu**:
- `existingElements`: `map[int]*entity.Element` - Ánh xạ số element sang con trỏ element
- `updatedNumbers`: `map[int]bool` - Theo dõi elements nào được xử lý từ request

**Độ phức tạp thời gian**: O(n) với n = số elements hiện tại

---

### 2. Vòng lặp xử lý chính (Dòng 644-780)

#### 2.1 Kiểm tra sự tồn tại của element

```go
for _, reqElem := range reqElements {
    existingElem, exists := existingElements[reqElem.Number]

    if exists {
        // Cập nhật element hiện có
    } else {
        // Tạo element mới
    }

    updatedNumbers[reqElem.Number] = true
}
```

**Logic**: Xử lý từng element từ request, hoặc cập nhật element hiện có hoặc tạo element mới.

---

## Logic xử lý chi tiết

### 3. Xử lý element hiện có (exists = true)

#### 3.1 Elements loại Picture

##### 3.1.1 Có khóa ảnh (không rỗng) (Dòng 652-685)

```go
if len(reqElem.PictureKeys) > 0 {
    // Chuyển đổi từ request sang định dạng entity
    reqPictureItems := make([]entity.PictureItem, len(reqElem.PictureKeys))
    for j, reqItem := range reqElem.PictureKeys {
        reqPictureItems[j] = entity.PictureItem{
            Key:   reqItem.Key,
            Order: reqItem.Order,
            Title: reqItem.Title,
        }
    }

    // Kiểm tra xem khóa có thay đổi không
    keysChanged := !u.pictureKeysEqual(existingElem.PictureKeys, reqPictureItems)

    if keysChanged {
        // Xóa có chọn lọc - chỉ xóa ảnh không có trong request mới
        keysToDelete := u.getKeysToDelete(existingElem.PictureKeys, reqPictureItems)
        for _, keyToDelete := range keysToDelete {
            if err := u.fileGateway.DeleteImage(ctx, keyToDelete); err != nil {
                log.Printf("failed to delete old picture: %v", err)
                continue
            }
        }
    }

    // Cập nhật với khóa ảnh mới
    existingElem.PictureKeys = make([]entity.PictureItem, len(reqElem.PictureKeys))
    for j, reqItem := range reqElem.PictureKeys {
        existingElem.PictureKeys[j] = entity.PictureItem{
            Key:   reqItem.Key,
            Order: reqItem.Order,
            Title: reqItem.Title,
        }
    }
}
```

**Thuật toán**: `getKeysToDelete`

```go
func (u *wikiUseCase) getKeysToDelete(currentKeys []entity.PictureItem, newKeys []entity.PictureItem) []string {
    currentKeyMap := make(map[string]bool)
    for _, item := range currentKeys {
        currentKeyMap[item.Key] = true
    }

    var keysToDelete []string
    for _, item := range newKeys {
        if currentKeyMap[item.Key] {
            delete(currentKeyMap, item.Key)
        }
    }

    for key := range currentKeyMap {
        keysToDelete = append(keysToDelete, key)
    }
    return keysToDelete
}
```

**Ví dụ thực tế**:
- **DB hiện tại**: `["anh1", "anh2", "anh3"]`
- **Request mới**: `["anh1", "anh4"]`
- **Khóa cần xóa**: `["anh2", "anh3"]` (không có trong request mới)
- **Kết quả**: Giữ `anh1`, xóa `anh2` & `anh3`, thêm `anh4`

##### 3.1.2 Mảng khóa ảnh rỗng (Dòng 686-698)

```go
} else {
    // Mảng rỗng - xóa tất cả ảnh hiện có
    if len(existingElem.PictureKeys) > 0 {
        for _, oldPictureItem := range existingElem.PictureKeys {
            if err := u.fileGateway.DeleteImage(ctx, oldPictureItem.Key); err != nil {
                log.Printf("failed to delete old picture: %v", err)
                continue
            }
        }
        // Xóa sạch khóa ảnh
        existingElem.PictureKeys = []entity.PictureItem{}
    }
}
```

**Ví dụ**: `["anh1", "anh2", "anh3"]` → `[]` → Xóa tất cả ảnh

#### 3.2 Các loại elements khác (Dòng 699-736)

```go
} else {
    // Xử lý các loại khác (text, image, file)
    newValue := reqElem.Value

    if newValue != nil && *newValue != "" {
        // Nếu giá trị thay đổi, xóa hình ảnh/file cũ
        if existingElem.Value != nil && *existingElem.Value != *newValue {
            if strings.EqualFold(existingElem.Type, "banner") {
                u.fileGateway.DeleteImage(ctx, *existingElem.Value)
            } else if strings.EqualFold(existingElem.Type, "large_picture") {
                u.fileGateway.DeleteImage(ctx, *existingElem.Value)
            } else if strings.EqualFold(existingElem.Type, "graphic") {
                u.fileGateway.DeleteImage(ctx, *existingElem.Value)
            } else if strings.EqualFold(existingElem.Type, "linked_in") {
                u.fileGateway.DeleteImage(ctx, *existingElem.Value)
            } else if strings.EqualFold(existingElem.Type, "file") {
                u.fileGateway.DeletePDF(ctx, *existingElem.Value)
            }
        }
    }

    existingElem.Value = reqElem.Value
}
```

**Các loại được hỗ trợ**: `banner`, `large_picture`, `graphic`, `linked_in`, `file`

### 4. Tạo element mới (exists = false) (Dòng 743-777)

```go
} else {
    // Element không tồn tại, thêm mới
    newElem := entity.Element{
        Number: reqElem.Number,
        Type:   reqElem.Type,
    }

    if strings.EqualFold(reqElem.Type, "picture") {
        if len(reqElem.PictureKeys) > 0 {
            newElem.PictureKeys = make([]entity.PictureItem, len(reqElem.PictureKeys))
            for j, reqItem := range reqElem.PictureKeys {
                newElem.PictureKeys[j] = entity.PictureItem{
                    Key:   reqItem.Key,
                    Order: reqItem.Order,
                    Title: reqItem.Title,
                }
            }
            // Đặt khóa đầu tiên làm giá trị chính để tương thích ngược
            newElem.Value = &reqElem.PictureKeys[0].Key
        }
    } else {
        if reqElem.Value != nil {
            newElem.Value = reqElem.Value
        }
    }

    if reqElem.VideoID != nil {
        newElem.VideoID = reqElem.VideoID
    }

    translation.Elements = append(translation.Elements, newElem)
}
```

**Lưu ý**: Với elements loại picture, khóa ảnh đầu tiên trở thành trường `Value` chính để tương thích ngược.

---

## Quản lý file

### 5. Giai đoạn dọn dẹp - Xóa elements đã xóa (Dòng 782-834)

```go
// Xóa elements không có trong request
var newElements []entity.Element
for _, existingElem := range translation.Elements {
    if updatedNumbers[existingElem.Number] {
        // Giữ elements đã được cập nhật
        newElements = append(newElements, existingElem)
    } else {
        // Xóa file của elements bị xóa
        if strings.EqualFold(existingElem.Type, "picture") {
            // Xóa tất cả ảnh trong mảng PictureKeys
            for _, pictureItem := range existingElem.PictureKeys {
                if pictureItem.Key != "" {
                    if err := u.fileGateway.DeleteImage(ctx, pictureItem.Key); err != nil {
                        return fmt.Errorf("failed to delete removed picture: %w", err)
                    }
                }
            }
        } else if existingElem.Value != nil && *existingElem.Value != "" {
            // Xóa file đơn dựa trên loại
            switch existingElem.Type {
            case "banner", "large_picture", "graphic", "linked_in":
                u.fileGateway.DeleteImage(ctx, *existingElem.Value)
            case "file":
                u.fileGateway.DeletePDF(ctx, *existingElem.Value)
            }
        }
        // Element bị xóa, không thêm vào newElements
    }
}

translation.Elements = newElements
```

**Logic**: Elements không có trong `updatedNumbers` sẽ bị xóa, và các file liên quan cũng bị xóa.

---

## Xử lý lỗi

### 6. Các trường hợp lỗi

1. **Lỗi xóa file**: Được ghi log nhưng không làm dừng quá trình (tiếp tục xử lý)
   ```go
   if err := u.fileGateway.DeleteImage(ctx, keyToDelete); err != nil {
       log.Printf("failed to delete old picture: %v", err)
       continue // Tiếp tục xóa các file khác
   }
   ```

2. **Lỗi xóa file quan trọng**: Trả về lỗi và dừng xử lý
   ```go
   if err := u.fileGateway.DeleteImage(ctx, pictureItem.Key); err != nil {
       return fmt.Errorf("failed to delete removed picture: %w", err)
   }
   ```

3. **Hủy bỏ Context**: Các thao tác tôn trọng việc hủy bỏ context để tắt máy gracefully

---

## Hiệu năng

### 7. Phân tích độ phức tạp thời gian

- **Chuẩn bị**: O(n) - n = số elements hiện tại
- **Xử lý**: O(m + k) - m = số elements request, k = thao tác khóa ảnh
- **Dọn dẹp**: O(n) - n = số elements hiện tại
- **Tổng cộng**: O(n + m + k)

### 8. Độ phức tạp không gian

- **Maps**: O(n) không gian thêm cho tra cứu elements
- **Mảng tạm thời**: O(k) cho xử lý khóa ảnh
- **Tổng cộng**: O(n + k)

### 9. Tối ưu hóa

1. **Tra cứu dựa trên Map**: O(1) truy cập elements thay vì tìm kiếm O(n)
2. **Xóa có chọn lọc**: Chỉ xóa ảnh đã thay đổi, không phải tất cả
3. **Xử lý theo batch**: Xử lý tất cả elements trong một lần duyệt

---

## Các trường hợp test

### 10. Các trường hợp test toàn diện

#### 10.1 Các trường hợp Picture Elements

**Trường hợp 1: Cập nhật với cùng khóa**
- Input: `["anh1", "anh2"]` → `["anh1", "anh2"]`
- Mong đợi: Không có thao tác file, chỉ cập nhật metadata

**Trường hợp 2: Thay thế một số ảnh**
- Input: `["anh1", "anh2", "anh3"]` → `["anh1", "anh4"]`
- Mong đợi: Xóa `anh2`, `anh3`; Giữ `anh1`; Thêm `anh4`

**Trường hợp 3: Xóa tất cả ảnh**
- Input: `["anh1", "anh2"]` → `[]`
- Mong đợi: Xóa `anh1`, `anh2`; Đặt PictureKeys = []

**Trường hợp 4: Thêm ảnh vào element rỗng**
- Input: `[]` → `["anh1", "anh2"]`
- Mong đợi: Thêm ảnh, không có thao tác xóa

#### 10.2 Các loại elements khác

**Trường hợp 5: Thay đổi ảnh banner**
- Input: `"banner_cu.jpg"` → `"banner_moi.jpg"`
- Mong đợi: Xóa `banner_cu.jpg`, cập nhật giá trị

**Trường hợp 6: Thay đổi file PDF**
- Input: `"tailieu_cu.pdf"` → `"tailieu_moi.pdf"`
- Mong đợi: Xóa `tailieu_cu.pdf`, cập nhật giá trị

#### 10.3 Quản lý Elements

**Trường hợp 7: Thêm element mới**
- Input: Element với số 5 (chưa tồn tại)
- Mong đợi: Tạo element mới, xử lý ảnh/files

**Trường hợp 8: Xóa element**
- Input: Element tồn tại trong DB nhưng không có trong request
- Mong đợi: Xóa element và tất cả file liên quan

**Trường hợp 9: Sắp xếp lại elements**
- Input: Elements [1,2,3] → [3,1,2]
- Mong đợi: Giữ tất cả elements, cập nhật thứ tự

---

## Trường hợp đặc biệt

### 11. Các trường hợp đặc biệt

1. **Request rỗng**: Tất cả elements hiện có bị xóa cùng với dọn dẹp file
2. **Số trùng lặp**: Element được xử lý sau cùng sẽ thắng (do map ghi đè)
3. **Giá trị nil**: Được xử lý an toàn, không crash
4. **Chuỗi rỗng**: Được coi là không có thay đổi giá trị
5. **Timeout Context**: Các thao tác có thể bị gián đoạn giữa chừng
6. **Lỗi hệ thống file**: Một số thao tác xóa có thể thất bại trong khi các thao tác khác thành công
7. **Mảng lớn**: Sử dụng bộ nhớ tỷ lệ thuận với số lượng ảnh
8. **Truy cập đồng thời**: Hàm không thread-safe, nên được gọi tuần tự

---

## Sơ đồ luồng

```
mergeElements()
├── Chuẩn bị Maps (existingElements, updatedNumbers)
├── Xử lý từng Element từ Request
│   ├── Element đã tồn tại?
│   │   ├── Có: Logic cập nhật
│   │   │   ├── Loại Picture?
│   │   │   │   ├── Có khóa?
│   │   │   │   │   ├── Có: So sánh & Xóa có chọn lọc
│   │   │   │   │   └── Không: Xóa tất cả ảnh
│   │   │   │   └── Loại khác: Kiểm tra thay đổi giá trị
│   │   │   └── Cập nhật trường Element
│   │   └── Không: Tạo Element mới
│   └── Đánh dấu đã cập nhật
├── Giai đoạn dọn dẹp
│   └── Xóa Elements chưa xử lý
│       └── Xóa Files liên quan
└── Cập nhật translation.Elements
```

---

## Ghi chú triển khai

### 12. Các quyết định thiết kế chính

1. **Xóa ảnh có chọn lọc**: Thay vì xóa tất cả ảnh khi có thay đổi, chỉ xóa những ảnh không có trong request mới
2. **Tương thích ngược**: Elements loại picture đặt khóa đầu tiên làm trường Value
3. **Khả năng phục hồi lỗi**: Lỗi xóa file không dừng quá trình xử lý
4. **Hiệu quả bộ nhớ**: Sử dụng maps cho tra cứu O(1) thay vì tìm kiếm tuyến tính
5. **Thao tác nguyên tử**: Hoặc tất cả thay đổi thành công hoặc duy trì trạng thái một phần

### 13. Cải tiến trong tương lai

1. **Hỗ trợ Transaction**: Bao bọc thao tác file trong transaction database
2. **Thao tác File theo batch**: Gom nhóm nhiều thao tác xóa file để hiệu quả hơn
3. **Lớp Validation**: Thêm validation đầu vào trước khi xử lý
4. **Audit Logging**: Theo dõi tất cả thay đổi để tuân thủ
5. **Retry Logic**: Thực hiện retry cho các thao tác file thất bại

---

*Tài liệu này bao gồm tất cả khía cạnh của hàm `mergeElements` theo implementation hiện tại. Khi có thay đổi logic, tài liệu này nên được cập nhật tương ứng.*
