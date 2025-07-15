
### `pouch.Inspect` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Inspect` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, belirtilen bir Docker nesnesinin (konteyner, imaj, ağ vb.) ayrıntılı yapılandırma ve durum bilgilerini almaktır. Bu işlem, arka planda `docker inspect` komutunu çalıştırarak gerçekleştirilir ve sonuç olarak genellikle JSON formatında bir metin döndürür.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// Inspect, belirtilen ID veya isme sahip Docker nesnesini denetler
// ve ayrıntılı bilgileri içeren bir JSON string'i döndürür.
func Inspect(id string) (string, error) {
	// "docker inspect [id]" komutunu oluştur.
	cmd := exec.Command("docker", "inspect", id)
	
	// Komutu çalıştır ve hem standart çıktıyı (stdout) hem de standart hatayı (stderr)
	// tek bir değişkende birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return "", fmt.Errorf("inspect hatası: %v\nçıktı: %s", err, out)
	}

	// Docker'ın çıktısının sonundaki olası boşlukları temizle ve string olarak döndür.
	return strings.TrimSpace(string(out)), nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "inspect", id)` satırı, Go'nun `os/exec` paketini kullanarak terminalde çalıştırılacak olan `docker inspect [id]` komutunu temsil eden bir nesne oluşturur. `id` parametresi, denetlenecek Docker nesnesinin kimliği veya adı olabilir.

2.  **Komutu Çalıştırma**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. `docker inspect` komutu, başarılı olduğunda nesneye ait tüm detayları JSON formatında standart çıktıya (stdout) yazar. Eğer belirtilen ID'ye sahip bir nesne bulunamazsa, hata mesajını standart hataya (stderr) yazar. Bu fonksiyon, her iki çıktıyı da tek bir byte dizisinde (`out`) birleştirerek, hem başarılı sonuçları hem de hata detaylarını yakalamayı kolaylaştırır.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun çalışması sırasında bir sorun olup olmadığını kontrol eder. Örneğin, var olmayan bir konteyner ID'si verilirse, `docker inspect` komutu hata verecek ve `err` değişkeni `nil` olmayacaktır. Bu durumda fonksiyon, orijinal hata (`err`) ve komutun ürettiği çıktıyı (`out`) birleştirerek zengin bir hata mesajı oluşturur ve geri döndürür. Bu, sorunun kaynağını anlamak için çok önemlidir.

4.  **Çıktının Temizlenmesi**:
    `strings.TrimSpace(string(out))` ifadesi, komut başarıyla çalıştıktan sonra elde edilen çıktıyı işler. `docker` komutları genellikle çıktılarının sonuna bir satır sonu karakteri (`\n`) ekler. `TrimSpace`, bu gibi baştaki ve sondaki tüm boşluk karakterlerini temizleyerek, saf ve işlenmesi daha kolay bir JSON metni elde edilmesini sağlar.

### Parametreler

*   `id (string)`: Denetlenecek Docker nesnesinin kimliği (tam veya kısaltılmış ID) veya adı. Bu bir konteyner, imaj, ağ, volume veya başka bir Docker nesnesi olabilir.

### Dönüş Değeri

*   `string`: İşlem başarılı olursa, denetlenen nesnenin tüm bilgilerini içeren, JSON formatında bir metin döndürür.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer belirtilen nesne bulunamazsa veya başka bir `docker` hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

Bu fonksiyonun düzgün çalışması için, Go programını çalıştıran sistemde **Docker'ın yüklü olması** ve `docker` komutunun sistemin `PATH` değişkeninde tanımlı olması gerekmektedir.

### Kullanım Örneği

Aşağıdaki örnekte, `my-app-container` adlı bir konteynerin IP adresini öğrenmek için `Inspect` fonksiyonunun nasıl kullanılacağı gösterilmiştir. `docker inspect` çıktısı JSON olduğu için, bu JSON'u Go'daki bir yapıya (`struct`) ayrıştırarak (unmarshal) istenen bilgiye kolayca erişilebilir.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

// Sadece ihtiyaç duyduğumuz alanları içeren basit bir struct
// docker inspect çıktısının tam yapısını temsil etmez.
type ContainerInspectInfo struct {
	NetworkSettings struct {
		IPAddress string `json:"IPAddress"`
	} `json:"NetworkSettings"`
}

func main() {
	containerName := "my-app-container"

	// pouch.Inspect'i kullanarak ham JSON çıktısını al
	jsonData, err := pouch.Inspect(containerName)
	if err != nil {
		log.Fatalf("Konteyner denetlenemedi: %v", err)
	}

	fmt.Printf("Ham JSON Çıktısı:\n%s\n\n", jsonData)

	// JSON çıktısını Go struct'ına ayrıştır
	var inspectInfo []ContainerInspectInfo
	if err := json.Unmarshal([]byte(jsonData), &inspectInfo); err != nil {
		log.Fatalf("JSON ayrıştırılamadı: %v", err)
	}

	// docker inspect her zaman bir dizi döndürdüğü için ilk elemanı alıyoruz.
	if len(inspectInfo) > 0 {
		ipAddress := inspectInfo[0].NetworkSettings.IPAddress
		fmt.Printf("'%s' konteynerinin IP adresi: %s\n", containerName, ipAddress)
	} else {
		log.Println("IP adresi bulunamadı.")
	}
}
```