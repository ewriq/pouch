

### `pouch.Start` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Start` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, durdurulmuş durumdaki bir Docker konteynerini başlatmaktır. Bu işlem, arka planda `docker start` komutunu çalıştırarak gerçekleştirilir ve işlemin yalnızca başarı veya başarısızlık durumunu bir `error` değeriyle bildirir.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// Start, belirtilen ID'ye sahip durdurulmuş bir Docker konteynerini başlatır.
func Start(id string) error {
	// "docker start [id]" komutunu oluştur.
	cmd := exec.Command("docker", "start", id)
	
	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	// docker start başarılı olduğunda konteynerin ID'sini çıktı olarak verir,
	// bu bilgi hata ayıklama için yakalanır.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return fmt.Errorf("start hatası: %v, çıktı: %s", err, string(out))
	}
	
	// İşlem başarılıysa, herhangi bir çıktı döndürmeden nil (hata yok) döndür.
	return nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "start", id)` satırı, terminalde `docker start [containerID]` komutunu çalıştırmak için bir komut nesnesi hazırlar.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. `docker start` komutu, bir konteyneri başarıyla başlattığında o konteynerin ID'sini standart çıktıya yazar. Bir hata durumunda ise (örneğin konteyner bulunamadığında) hata mesajını standart hataya yazar. `CombinedOutput` fonksiyonu, her iki çıktıyı da tek bir `out` değişkeninde toplayarak, özellikle hata durumlarında Docker'ın döndürdüğü mesajı görmeyi sağlar.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun sıfırdan farklı bir çıkış koduyla sonlanıp sonlanmadığını kontrol eder. Eğer belirtilen ID'ye sahip bir konteyner mevcut değilse veya Docker servisi (daemon) çalışmıyorsa, `docker start` komutu hata verir. Fonksiyon bu hatayı yakalar ve `fmt.Errorf` kullanarak hem Go'nun genel hata bilgisini (`err`) hem de Docker'ın ürettiği spesifik çıktıyı (`out`) içeren detaylı bir hata mesajı oluşturur.

4.  **Başarılı Durum**:
    Komut başarıyla çalışırsa, `err` değişkeni `nil` olur. Bu fonksiyon, başarılı olduğunda komutun ürettiği çıktıyı (başlatılan konteynerin ID'si) döndürmez. Bunun yerine, sadece işlemin başarılı olduğunu belirtmek için `nil` değeri döndürür. Bu, "sadece yap ve sonucunu bildir" (fire-and-report) tarzı basit bir tasarım tercihidir.

### Parametreler

*   `id (string)`: Başlatılacak olan durdurulmuş Docker konteynerinin kimliği (ID) veya adı.

### Dönüş Değeri

*   `error`: Fonksiyon bir `error` değeri döndürür.
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner bulunamazsa, zaten çalışıyorsa (bu bir hata değildir ama yeniden başlatmaz) veya başka bir Docker hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

*   Bu fonksiyonun çalışabilmesi için, kodu çalıştıran sistemde **Docker'ın yüklü olması** ve `docker` komutunun sistemin `PATH` değişkeni üzerinden erişilebilir olması gerekmektedir.

### Kullanım Örneği

Sistemde `my-database-container` adında durdurulmuş bir konteyneri başlatmak için `Start` fonksiyonunun nasıl kullanılacağını gösteren bir örnek:

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerToStart := "my-database-container"

	log.Printf("'%s' konteyneri başlatılıyor...", containerToStart)

	// pouch.Start fonksiyonunu çağır
	err := pouch.Start(containerToStart)
	if err != nil {
		log.Fatalf("Konteyner başlatılamadı: %v", err)
	}

	log.Printf("'%s' konteyneri başarıyla başlatıldı!", containerToStart)
}
```