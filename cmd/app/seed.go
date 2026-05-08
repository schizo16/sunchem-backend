package app

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedDB(db *gorm.DB) {
	var count int64
	db.Raw("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 {
		fmt.Println("[Seed] DB already has data, skipping")
		return
	}

	fmt.Println("[Seed] Populating products...")
	products := []struct {
		Slug, Name, Desc, Image, Category, Highlights string
	}{
		{Slug: "nhua-nhu-tuong-sunperse-c-50", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-50", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "SUNPERSE C-50 là nhựa nhũ tương trên cơ sở copolymer styrene acrylic, có các hạt nhũ với kích thước siêu nhỏ và đồng nhất. Sunperse C-50 là chất tạo màng tương thích tốt với các loại bột màu và chất độn có khả năng kết dính rất cao.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/378e720bc90705595c16.jpg?v=1659003164017", Highlights: "Chất tạo màng,Chịu chà rửa tốt,Dẻo dai,Kết dính cao"},
		{Slug: "nhua-nhu-tuong-sunperse-c-68", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-68", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "SUNPERSE C-68 là nhựa nhũ tương trên cơ sở copolymer styrene acrylic, được dùng làm chất tạo màng trong các hệ chống thấm trộn xi măng. Sunperse C-68 có khả năng kháng nước rất cao cùng độ mềm dẻo cực tốt.", Image: "https://bizweb.dktcdn.net/100/229/353/products/378e720bc90705595c16.jpg", Highlights: "Chống thấm,Kháng nước cao,Mềm dẻo,Hệ xi măng"},
		{Slug: "nhua-nhu-tuong-sunperse-c-77", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-77", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "SUNPERSE C-77 là nhựa nhũ tương trên cơ sở copolymer styrene acrylic, có các hạt nhũ với kích thước siêu nhỏ và đồng nhất. Sản phẩm có khả năng kết dính bột màu vượt trội.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/products/378e720bc90705595c16.jpg", Highlights: "Kết dính bột màu,Không dính tay,Tương thích tốt"},
		{Slug: "titanium-dioxit-rutile-r668", Name: "RUTILE TITAN DIOXIT R668", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "TiO2 R668 có hàm lượng TiO2 93%, ánh vàng, hạt titan được bọc nhiều lớp, độ phủ, độ bền thời tiết tốt, thường được dùng cho các dòng sơn nội thất và ngoại thất kinh tế.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/files/12344-989f24af-27d7-483d-b572-292325b4bb63.jpg", Highlights: "TiO2 93%,Ánh vàng,Sơn nội/ngoại thất,Độ bền thời tiết"},
		{Slug: "phan-tan-cao-cap-pidicryl-120v", Name: "PHÂN TÁN CAO CẤP PIDICRYL 120V", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "Pidicryl 120v là một loại phụ gia phân tán cao cấp dùng cho hệ sơn và mực in gốc nước. Có khả năng phân tán và ổn định bột màu vô cơ và hữu cơ, đặc biệt là bột màu đen carbon.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/products/pidicryl-120v.jpg", Highlights: "Phân tán cao cấp,Ổn định bột màu đen,Cải thiện lưu biến,Sơn & mực in"},
		{Slug: "mau-paste-goc-nuoc", Name: "MÀU PHA MÁY SMARTIN", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "Màu paste các màu hữu cơ hoặc vô cơ đã được nghiền mịn sẵn thành dạng sệt phân tán trong nước. Khi pha màu sơn, chỉ cần cho màu vào base trắng và điều chỉnh tỷ lệ phù hợp.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/files/z3599983022221-a75b3515222155716902cb041a010f7f.jpg", Highlights: "Màu paste,Pha máy tự động,Độ bền cao,Phân phối độc quyền"},

		{Slug: "cabelec-xs6639c", Name: "CABELEC XS6639C", Category: "HÓA CHẤT NGÀNH NHỰA", Desc: "CABELEC XS6639C là một loại hạt chống tĩnh điện cao cấp của CABOT, chứa hàm lượng cao hạt carbon ở dạng ống nano, được thiết kế để loại bỏ hiện tượng tĩnh điện trên các sản phẩm nhựa PE, PP.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/products/350928849-804434641263328-5543458960280316847-n.jpg", Highlights: "Chống tĩnh điện,Công nghệ nano,CABOT USA,PE & PP"},
		{Slug: "plasblak-pe2705", Name: "PLASBLAK PE2705", Category: "HÓA CHẤT NGÀNH NHỰA", Desc: "PLASBLAK PE2705 là hạt masterbatch màu đen, chứa tới 50% hạt carbon trên nhựa nền PE. Đặc biệt PE2705 có thể dùng để nhuộm màu đen cho các loại nhựa ABS, SAN, PS.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/products/350928849-804434641263328-5543458960280316847-n-b8458c02-fe9a-4011-ae19-f470e0d23991.jpg", Highlights: "Carbon 50%,Đa nền nhựa,CABOT USA,Màu đen sâu"},
		{Slug: "titanium-dioxit-rutile-r902", Name: "RUTILE TITAN DIOXIT R902", Category: "HÓA CHẤT NGÀNH MỰC IN", Desc: "TiO2 R902+ có Hàm lượng TiO2 93%, ánh trung tính, được sử dụng nhiều trong các dòng sơn gỗ, sơn dung môi.", Image: "https://bizweb.dktcdn.net/thumb/large/100/229/353/files/12344-989f24af-27d7-483d-b572-292325b4bb63.jpg", Highlights: "TiO2 93%,Ánh trung tính,Sơn gỗ & dung môi,Chemours"},
		{Slug: "titanium-dioxit-rutile-r706-1", Name: "RUTILE TITAN DIOXIT R706", Category: "HÓA CHẤT NGÀNH MỰC IN", Desc: "TiO2 R706 có Hàm lượng TiO2 93%, ánh xanh, hạt titan được xử lý bọc nhiều lớp cho độ phủ, độ bền thời tiết tốt, được sử dụng nhiều trong các dòng sơn nước ngoại thất và nội thất cao cấp.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/706.jpg", Highlights: "TiO2 93%,Ánh xanh,Sơn nước cao cấp,Chemours"},

		{Slug: "hat-mau-den-cabot-plasblak-un2005", Name: "Hạt màu đen Cabot PLASBLAK UN2005", Category: "HÓA CHẤT NGÀNH NHỰA", Desc: "Hạt màu đen Cabot PLASBLAK UN2005 có hàm lượng carbon 50%, là hạt black masterbatch cao cấp chuyên cho những ứng dụng cần màu đen sâu và bóng, ánh xanh, phù hợp với đa dạng nền nhựa như PE, PP, ABS, PS, PC.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/350928849-804434641263328-5543458960280316847-n-209d21c5-8602-483b-8b1f-163b95c9287d.jpg?v=1703752651347", Highlights: "Carbon 50%,Đa nền nhựa PE/PP/ABS,Màu đen sâu ánh xanh"},
		{Slug: "hat-mau-den-cabot-plasblak-xp-6603a", Name: "Hạt màu đen Cabot PLASBLAK XP 6603A", Category: "HÓA CHẤT NGÀNH NHỰA", Desc: "Hạt màu đen Cabot PLASBLAK XP 6603A, hàm lượng carbon 45% cho màu đen sâu, bóng, ánh xanh, ứng dụng trong sản xuất các sản phẩm nhựa thổi film, đùn ống, compound, ép phun.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/350928849-804434641263328-5543458960280316847-n-d39919d7-a56a-4c48-b004-57b5d6563541.jpg?v=1703752229597", Highlights: "Carbon 45%,Thổi film,Đùn ống,Ép phun"},
		{Slug: "hat-mau-den-cabot-plasblak-pe-2705", Name: "Hạt màu đen Cabot PLASBLAK PE 2705", Category: "HÓA CHẤT NGÀNH NHỰA", Desc: "Hạt màu đen Cabot PLASBLAK PE 2705, hàm lượng carbon 50% cho màu đen sâu, bóng, ánh xanh, ứng dụng trong sản xuất các sản phẩm nhựa ép phun, cán màng, thổi can, thổi film, compound.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/350928849-804434641263328-5543458960280316847-n-b8458c02-fe9a-4011-ae19-f470e0d23991.jpg?v=1703751863247", Highlights: "Carbon 50%,Ép phun,Cán màng,Thổi film"},
		{Slug: "hat-mau-den-hat-chong-tinh-dien-cabot-cabot-black-masterbatch-conductive", Name: "Hạt màu đen Cabot PLASBLAK PE 2718", Category: "HÓA CHẤT NGÀNH NHỰA", Desc: "Hạt màu đen Cabot PLASBLAK PE 2718, hàm lượng carbon 50% cho màu đen sâu, bóng, ánh xanh, ứng dụng trong sản xuất các sản phẩm nhựa ép phun, cán màng, thổi can, thổi film, compound.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/350928849-804434641263328-5543458960280316847-n.jpg?v=1696325329633", Highlights: "Carbon 50%,Chống tĩnh điện,CABOT,Ép phun"},
		{Slug: "titanium-dioxit-rutile-r706", Name: "TITANIUM DIOXIT RUTILE R706", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "TiO2 R706 có Hàm lượng TiO2 93%, ánh xanh, hạt titan được xử lý bọc nhiều lớp cho độ phủ, độ bền thời tiết tốt, được sử dụng nhiều trong các dòng sơn nước ngoại thất và nội thất cao cấp.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/706.jpg?v=1703753349267", Highlights: "TiO2 93%,Ánh xanh,Chemours USA,Sơn nước ngoại thất"},
		{Slug: "chat-tao-mang-binder", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-50", Category: "HÓA CHẤT NGÀNH SƠN", Desc: "SUNPERSE C-50 là nhựa nhũ tương trên cơ sở copolymer styrene acrylic, có các hạt nhũ với kích thước siêu nhỏ và đồng nhất. Chất tạo màng thích hợp cho sơn nước.", Image: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/products/378e720bc90705595c16.jpg?v=1659003164017", Highlights: "Chất tạo màng,Kết dính cao,Chống thấm,Sơn nước"},
	}

	for _, p := range products {
		db.Exec(`INSERT OR IGNORE INTO products (slug, name, short_description, image, category, highlights, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)`,
			p.Slug, p.Name, p.Desc, p.Image, p.Category, p.Highlights, time.Now(), time.Now())
	}
	fmt.Printf("[Seed] Inserted %d products\n", len(products))

	fmt.Println("[Seed] Populating blog posts...")
	blogs := []struct {
		Slug, Title, Summary, Content, Thumbnail, Date string
	}{
		{
			Slug: "hoi-thao-melodies-of-colors-giai-dieu-cua-sac-mau",
			Title: `HỘI THẢO "MELODIES OF COLORS - GIAI ĐIỆU CỦA SẮC MÀU"`,
			Summary: "Công ty TNHH Sunchem xin chân thành cảm ơn Quý khách hàng đã tham dự hội thảo Melodies of Colors do chúng tôi cùng các Nhà cung cấp Huiyun, Soujanya và Italtinto tổ chức.",
			Content: `<p>Công ty TNHH Sunchem xin chân thành cảm ơn Quý khách hàng đã dành thời gian quý báu để tham dự hội thảo "Melodies of Colors - Giai điệu của Sắc màu" do chúng tôi cùng các Nhà cung cấp Huiyun, Soujanya và Italtinto tổ chức.</p><p>Hội thảo lần này không chỉ nhằm mục đích giới thiệu những sản phẩm và công nghệ tiên tiến trong lĩnh vực hóa chất ngành sơn, mà còn là dịp để chúng ta gặp gỡ, trao đổi và cùng nhau phát triển.</p><p>Một lần nữa, chúng tôi trân trọng cảm ơn và rất mong sẽ tiếp tục nhận được sự đồng hành, hợp tác của Quý khách trong thời gian tới.</p>`,
			Thumbnail: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/articles/son07789.jpg?v=1722833711910",
			Date: "2024-08-05",
		},
		{
			Slug: "thong-bao-lich-nghi-le-30-4-1-5",
			Title: "THÔNG BÁO LỊCH NGHỈ LỄ 30/4-1/5",
			Summary: "Công ty Sunchem xin gửi tới Quý Công ty lịch nghỉ lễ Ngày Giải phóng miền Nam 30/4 và Quốc tế Lao động 1/5.",
			Content: `<p>Kính gửi: Quý khách hàng, Quý đối tác</p><p>Công ty Sunchem xin gửi tới Quý Công ty lịch nghỉ lễ Ngày Giải phóng miền Nam 30/4 và Quốc tế Lao động 1/5 của chúng tôi như sau:</p><p><strong>Ngày nghỉ: bắt đầu từ Thứ 7 (27/04/2024) đến hết Thứ 4 (01/05/2024)</strong></p><p><strong>Ngày làm việc trở lại: Thứ 5 (02/05/2024)</strong></p>`,
			Thumbnail: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/articles/beige-and-red-holiday-closure-notice-instagram-post.jpg?v=1713930599577",
			Date: "2024-04-24",
		},
		{
			Slug: "thong-bao-lich-nghi-tet-nguyen-dan-2024",
			Title: "THÔNG BÁO LỊCH NGHỈ TẾT NGUYÊN ĐÁN 2024",
			Summary: "CÔNG TY SUNCHEM xin chân thành cảm ơn sự ủng hộ của quý khách hàng trong năm vừa qua. Chúc Quý khách một năm mới 2024 An khang thịnh vượng!",
			Content: `<p>CÔNG TY SUNCHEM xin chân thành cảm ơn sự ủng hộ của quý khách hàng trong năm vừa qua. Chúc Quý khách một năm mới 2024 Khởi sắc, Dồi dào sức khỏe, An khang thịnh vượng, Vạn sự như ý!</p><p><strong>Thời gian nghỉ: Bắt đầu từ Thứ 4, ngày 07/02/2024 đến hết Thứ 4, ngày 14/02/2024.</strong></p><p><strong>Thời gian làm việc trở lại: Thứ 5, Ngày 15/02/2024.</strong></p>`,
			Thumbnail: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/articles/happy-new-year-2024.png?v=1706522833883",
			Date: "2024-01-29",
		},
		{
			Slug: "thong-bao-lich-nghi-le-tet-duong-lich-2024-1",
			Title: "THÔNG BÁO LỊCH NGHỈ LỄ TẾT DƯƠNG LỊCH 2024",
			Summary: "Công ty TNHH Sunchem xin chân thành cảm ơn sự hợp tác của Quý khách hàng trong thời gian qua. Nhân dịp năm mới, chúc Quý khách một năm 2024 thành công, thịnh vượng.",
			Content: `<p>Công ty TNHH Sunchem xin chân thành cảm ơn sự hợp tác của Quý khách hàng, Quý nhà cung cấp trong thời gian qua.</p><p>Nhân dịp năm mới, công ty Sunchem xin gửi lời chúc tới Quý khách hàng, Quý nhà cung cấp cùng gia đình một kỳ nghỉ lễ vui vẻ, hạnh phúc, một năm 2024 thành công, thịnh vượng.</p><p><strong>Bắt đầu từ Thứ Bảy ngày 30/12/2023 đến hết Thứ Hai ngày 01/01/2024.</strong></p><p>Thứ Ba ngày 02/01/2024 công ty chúng tôi trở lại làm việc bình thường</p>`,
			Thumbnail: "https://bizweb.dktcdn.net/thumb/grande/100/229/353/articles/thiet-ke-chua-co-ten.png?v=1703758583013",
			Date: "2023-12-28",
		},
	}

	for _, b := range blogs {
		pubTime, _ := time.Parse("2006-01-02", b.Date)
		db.Exec(`INSERT OR IGNORE INTO blog_posts (title, slug, summary, content, thumbnail, category, status, views, published_at, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			b.Title, b.Slug, b.Summary, b.Content, b.Thumbnail, "Tin tức", "published", 0, pubTime, time.Now(), time.Now())
	}
	fmt.Printf("[Seed] Inserted %d blog posts\n", len(blogs))
}
