### `pouch.Pull` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Pull` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, belirtilen bir Docker imajını bir imaj kayıt defterinden (genellikle Docker Hub) yerel sisteme indirmektir. Bu işlem, arka planda `docker pull` komutunu çalıştırarak gerçekleştirilir ve işlemin başarı durumunu bir `error` değeriyle bildirir.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// Pull, belirtilen Docker imajını uzak bir kayıt defterinden çeker (indirir).
func Pull(image string) error {
	// "docker pull [image]" komutunu oluştur.
	cmd := exec.Command("docker", "pull", image)

	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	// Bu, hem ilerleme çubuğu gibi çıktıları hem de olası hataları yakalamayı sağlar.
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem Go hatasını hem de Docker'ın çıktısını içeren
		// açıklayıcı bir hata mesajı döndür.
		return fmt.Errorf("pull hatası: %v\n%s", err, string(output))
	}

	// İşlem başarılıysa nil (hata yok) döndür.
	return nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "pull", image)` satırı, terminalde çalıştırılacak olan `docker pull [imaj_adı]` komutunu programatik olarak oluşturur. `image` parametresi, "ubuntu:latest" veya "redis" gibi bir imaj adı olabilir.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. `docker pull` komutu çalışırken indirme ilerlemesi gibi bilgileri standart çıktıya, hata mesajlarını ise standart hataya yazar. `CombinedOutput`, bu iki akışı da tek bir byte dizisinde (`output`) toplar. Bu, hata durumunda Docker'ın döndürdüğü spesifik hata mesajını (örneğin, "repository not found") yakalamak için çok etkilidir.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun başarılı bir şekilde tamamlanıp tamamlanmadığını kontrol eder. Eğer imaj bulunamazsa, internet bağlantısı yoksa veya Docker kayıt defterinde bir sorun varsa, `docker pull` komutu sıfırdan farklı bir çıkış koduyla sonlanır ve `err` değişkeni `nil` olmaz. Bu durumda fonksiyon, `fmt.Errorf` ile hem Go seviyesindeki genel hatayı (`err`) hem de Docker CLI'den gelen detaylı çıktıyı (`output`) içeren zengin bir hata mesajı oluşturarak geri döner.

4.  **Başarılı Durum**:
    Eğer imaj başarıyla indirilirse, `docker pull` komutu 0 çıkış koduyla tamamlanır ve `err` değeri `nil` olur. Fonksiyon bu durumda `nil` döndürerek işlemin başarılı olduğunu belirtir.

### Parametreler

*   `image (string)`: Çekilecek (indirilecek) Docker imajının adı ve etiketi (tag). Örneğin: `"alpine:latest"`, `"postgres:14"`.

### Dönüş Değeri

*   `error`: Fonksiyon bir `error` değeri döndürür.
    *   İmaj indirme işlemi başarılı olursa bu değer `nil` olur.
    *   İmaj bulunamazsa, ağ bağlantısı sorunları yaşanırsa veya başka bir Docker hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

*   Bu fonksiyonun çalışabilmesi için, Go kodunun çalıştığı sistemde **Docker'ın yüklü olması** ve `docker` komutunun sistemin `PATH` değişkeninde bulunması gerekir.
*   Uzak bir kayıt defterinden imaj çekileceği için aktif bir **internet bağlantısı** gereklidir.

### Kullanım Örneği

Sisteme `alpine` imajının en son sürümünü çekmek için `Pull` fonksiyonunun nasıl kullanılacağını gösteren bir örnek:

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	imageToPull := "alpine:latest"

	log.Printf("'%s' imajı çekiliyor...", imageToPull)

	// pouch.Pull fonksiyonunu çağır
	err := pouch.Pull(imageToPull)
	if err != nil {
		log.Fatalf("İmaj çekme işlemi başarısız oldu: %v", err)
	}

	log.Printf("'%s' imajı başarıyla çekildi!", imageToPull)
}
```