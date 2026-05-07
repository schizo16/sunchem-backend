package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	sunchemDB "sunchem-backend/internal/common/db"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("G:/code/ssss/sunchem-backend/sunchem.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := sunchemDB.AutoMigrate(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Crawl] DB migrated OK")
}

type ProductEntry struct {
	Slug     string
	Name     string
	Category string
	ImageURL string
}

type BlogEntry struct {
	Slug  string
	Title string
	Date  string
}

func slugify(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func fetchHTML(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func stripTags(html string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return strings.TrimSpace(re.ReplaceAllString(html, ""))
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

func crawlProduct(entry ProductEntry) {
	url := "https://sunchemvn.com/" + entry.Slug
	fmt.Printf("  Fetching product: %s\n", entry.Name)
	html, err := fetchHTML(url)
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
		return
	}

	image := ""
	imgRe := regexp.MustCompile(`<meta property="og:image" content="([^"]+)"`)
	imgMatches := imgRe.FindStringSubmatch(html)
	if len(imgMatches) > 1 {
		image = imgMatches[1]
	}
	if image == "" {
		imgRe2 := regexp.MustCompile(`(https?://bizweb\.dktcdn\.net/100/229/353/products/[^"'\s]+)`)
		imgMatches2 := imgRe2.FindAllString(html, 1)
		if len(imgMatches2) > 0 {
			image = imgMatches2[0]
		}
	}

	desc := ""
	descRe := regexp.MustCompile(`<meta property="og:description" content="([^"]+)"`)
	descMatches := descRe.FindStringSubmatch(html)
	if len(descMatches) > 1 {
		desc = strings.TrimSpace(descMatches[1])
	}
	if desc == "" {
		descRe2 := regexp.MustCompile(`<meta name="description" content="([^"]+)"`)
		descMatches2 := descRe2.FindStringSubmatch(html)
		if len(descMatches2) > 1 {
			desc = strings.TrimSpace(descMatches2[1])
		}
	}

	slug := entry.Slug
	if slug == "" {
		slug = slugify(entry.Name)
	}

	fmt.Printf("    Image: %s\n", truncate(image, 60))
	fmt.Printf("    Desc: %s\n", truncate(desc, 80))

	var count int64
	db.Raw("SELECT COUNT(*) FROM products WHERE slug = ?", slug).Scan(&count)
	if count > 0 {
		db.Exec(`UPDATE products SET name=?, short_description=?, image=?, category=?, highlights=?, updated_at=? WHERE slug=?`,
			entry.Name, desc, image, entry.Category, entry.Category, time.Now(), slug)
		fmt.Printf("    UPDATED\n")
	} else {
		db.Exec(`INSERT INTO products (slug, name, short_description, image, category, highlights, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)`,
			slug, entry.Name, desc, image, entry.Category, entry.Category, time.Now(), time.Now())
		fmt.Printf("    INSERTED\n")
	}
}

func crawlBlog(entry BlogEntry) {
	url := "https://sunchemvn.com/" + entry.Slug
	fmt.Printf("  Fetching blog: %s\n", entry.Title)
	html, err := fetchHTML(url)
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
		return
	}

	content := ""
	rteRe := regexp.MustCompile(`(?s)<div class="rte">(.*?)</div>`)
	rteMatch := rteRe.FindStringSubmatch(html)
	if len(rteMatch) > 1 {
		content = strings.TrimSpace(rteMatch[1])
	}
	cleanContent := stripTags(content)

	thumbnail := ""
	imgRe := regexp.MustCompile(`<meta property="og:image" content="([^"]+)"`)
	imgMatch := imgRe.FindStringSubmatch(html)
	if len(imgMatch) > 1 {
		thumbnail = imgMatch[1]
	}
	if thumbnail == "" {
		imgRe2 := regexp.MustCompile(`(https?://bizweb\.dktcdn\.net/100/229/353/articles/[^"'\s]+)`)
		imgMatch2 := imgRe2.FindAllString(html, 1)
		if len(imgMatch2) > 0 {
			thumbnail = imgMatch2[0]
		}
	}

	date := entry.Date
	if date == "" {
		dateRe := regexp.MustCompile(`lúc (\d{2}/\d{2}/\d{4})`)
		dateMatch := dateRe.FindStringSubmatch(html)
		if len(dateMatch) > 1 {
			date = dateMatch[1]
		}
	}

	summary := cleanContent
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}

	pubTime, _ := time.Parse("02/01/2006", date)

	slug := entry.Slug
	if slug == "" {
		slug = slugify(entry.Title)
	}

	fmt.Printf("    Thumbnail: %s\n", truncate(thumbnail, 60))
	fmt.Printf("    Content len: %d\n", len(content))
	fmt.Printf("    Date: %s\n", date)

	var count int64
	db.Raw("SELECT COUNT(*) FROM blog_posts WHERE slug = ?", slug).Scan(&count)
	if count > 0 {
		db.Exec(`UPDATE blog_posts SET title=?, summary=?, content=?, thumbnail=?, category=?, published_at=?, updated_at=? WHERE slug=?`,
			entry.Title, summary, content, thumbnail, "Tin tức", pubTime, time.Now(), slug)
		fmt.Printf("    UPDATED\n")
	} else {
		db.Exec(`INSERT INTO blog_posts (title, slug, summary, content, thumbnail, category, status, views, published_at, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			entry.Title, slug, summary, content, thumbnail, "Tin tức", "published", 0, pubTime, time.Now(), time.Now())
		fmt.Printf("    INSERTED\n")
	}
}

func main() {
	products := []ProductEntry{
		{Slug: "nhua-nhu-tuong-sunperse-c-50", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-50", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "nhua-nhu-tuong-sunperse-c-68", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-68", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "nhua-nhu-tuong-sunperse-c-77", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-77", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "titanium-dioxit-rutile-r668", Name: "RUTILE TITAN DIOXIT R668", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "phan-tan-cao-cap-pidicryl-120v", Name: "PHÂN TÁN CAO CẤP PIDICRYL 120V", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "mau-paste-goc-nuoc", Name: "MÀU PHA MÁY SMARTIN", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "cabelec-xs6639c", Name: "CABELEC XS6639C", Category: "HÓA CHẤT NGÀNH NHỰA"},
		{Slug: "plasblak-pe2705", Name: "PLASBLAK PE2705", Category: "HÓA CHẤT NGÀNH NHỰA"},
		{Slug: "titanium-dioxit-rutile-r902", Name: "RUTILE TITAN DIOXIT R902", Category: "HÓA CHẤT NGÀNH MỰC IN"},
		{Slug: "titanium-dioxit-rutile-r706-1", Name: "RUTILE TITAN DIOXIT R706", Category: "HÓA CHẤT NGÀNH MỰC IN"},
		{Slug: "hat-mau-den-cabot-plasblak-un2005", Name: "Hạt màu đen Cabot PLASBLAK UN2005", Category: "HÓA CHẤT NGÀNH NHỰA"},
		{Slug: "hat-mau-den-cabot-plasblak-xp-6603a", Name: "Hạt màu đen Cabot PLASBLAK XP 6603A", Category: "HÓA CHẤT NGÀNH NHỰA"},
		{Slug: "hat-mau-den-cabot-plasblak-pe-2705", Name: "Hạt màu đen Cabot PLASBLAK PE 2705", Category: "HÓA CHẤT NGÀNH NHỰA"},
		{Slug: "hat-mau-den-hat-chong-tinh-dien-cabot-cabot-black-masterbatch-conductive", Name: "Hạt màu đen Cabot PLASBLAK PE 2718", Category: "HÓA CHẤT NGÀNH NHỰA"},
		{Slug: "titanium-dioxit-rutile-r706", Name: "TITANIUM DIOXIT RUTILE R706", Category: "HÓA CHẤT NGÀNH SƠN"},
		{Slug: "chat-tao-mang-binder", Name: "NHỰA NHŨ TƯƠNG SUNPERSE C-50", Category: "HÓA CHẤT NGÀNH SƠN"},
	}

	blogs := []BlogEntry{
		{Slug: "hoi-thao-melodies-of-colors-giai-dieu-cua-sac-mau", Title: "HỘI THẢO \"MELODIES OF COLORS\"", Date: "05/08/2024"},
		{Slug: "thong-bao-lich-nghi-le-30-4-1-5", Title: "THÔNG BÁO LỊCH NGHỈ LỄ 30/4-1/5", Date: "24/04/2024"},
		{Slug: "thong-bao-lich-nghi-tet-nguyen-dan-2024", Title: "THÔNG BÁO LỊCH NGHỈ TẾT NGUYÊN ĐÁN 2024", Date: "29/01/2024"},
		{Slug: "thong-bao-lich-nghi-le-tet-duong-lich-2024-1", Title: "THÔNG BÁO LỊCH NGHỈ LỄ TẾT DƯƠNG LỊCH 2024", Date: "28/12/2023"},
	}

	fmt.Println("=== CRAWLING PRODUCTS ===")
	for _, p := range products {
		crawlProduct(p)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n=== CRAWLING BLOG POSTS ===")
	for _, b := range blogs {
		crawlBlog(b)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n=== DONE ===")
	db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	var productCount, blogCount int64
	db.Raw("SELECT COUNT(*) FROM products").Scan(&productCount)
	db.Raw("SELECT COUNT(*) FROM blog_posts").Scan(&blogCount)
	fmt.Printf("Products in DB: %d\n", productCount)
	fmt.Printf("Blog posts in DB: %d\n", blogCount)
}
