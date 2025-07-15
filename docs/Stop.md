
### `pouch.Stop` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Stop` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, çalışan bir Docker konteynerini güvenli bir şekilde durdurmaktır. Bu işlem, arka planda `docker stop` komutunu çalıştırarak gerçekleştirilir ve işlemin başarı veya başarısızlık durumunu bir `error` değeriyle bildirir.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// Stop, belirtilen ID'ye sahip bir Docker konteynerini durdurur.
func Stop(id string) error {
	// "docker stop [id]" komutunu oluştur.
	cmd := exec.Command("docker", "stop", id)

	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	// docker stop başarılı olduğunda konteynerin ID'sini çıktı olarak verir,
	// bu bilgi hata ayıklama için yakalanır.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		// Not: Orijinal koddaki "start error" muhtemelen bir yazım hatasıydı, "stop error" olmalı.
		return fmt.Errorf("stop hatası: %v, çıktı: %s", err, string(out))
	}
	
	// İşlem başarılıysa nil (hata yok) döndür.
	return nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "stop", id)` satırı, terminalde `docker stop [containerID]` komutunu çalıştırmak için bir komut nesnesi hazırlar. `docker stop` komutu, konteynere varsayılan olarak bir `SIGTERM` sinyali gönderir ve belirli bir süre içinde (genellikle 10 saniye) kapanmasını bekler. Kapanmazsa bir `SIGKILL` sinyali göndererek zorla sonlandırır.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. Bir konteyner başarıyla durdurulduğunda, `docker stop` komutu durdurulan konteynerin ID'sini standart çıktıya yazar. Bir hata durumunda ise (örneğin konteyner bulunamadığında) hata mesajını standart hataya yazar. `CombinedOutput` fonksiyonu, her iki çıktıyı da tek bir `out` değişkeninde toplayarak, özellikle hata durumlarında Docker'ın döndürdüğü mesajı görmeyi sağlar.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun sıfırdan farklı bir çıkış koduyla sonlanıp sonlanmadığını kontrol eder. Eğer belirtilen ID'ye sahip bir konteyner mevcut değilse veya Docker servisi (daemon) çalışmıyorsa, `docker stop` komutu hata verir. Fonksiyon bu hatayı yakalar ve `fmt.Errorf` kullanarak hem Go'nun genel hata bilgisini (`err`) hem de Docker'ın ürettiği spesifik çıktıyı (`out`) içeren detaylı bir hata mesajı oluşturur.

4.  **Başarılı Durum**:
    Komut başarıyla çalışırsa, `err` değişkeni `nil` olur. Fonksiyon, işlemin başarılı olduğunu belirtmek için `nil` değeri döndürür. Başarılı durumda konteynerin ID'sini döndürmez, çünkü asıl amaç işlemi gerçekleştirmektir.

### Parametreler

*   `id (string)`: Durdurulacak olan çalışan Docker konteynerinin kimliği (ID) veya adı.

### Dönüş Değeri

*   `error`: Fonksiyon bir `error` değeri döndürür.
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner bulunamazsa, zaten durdurulmuşsa veya başka bir Docker hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

*   Bu fonksiyonun çalışabilmesi için, kodu yürüten sistemde **Docker'ın yüklü olması** ve `docker` komutunun sistemin `PATH` değişkeni üzerinden erişilebilir olması gerekmektedir.

### Kullanım Örneği

Sistemde `my-nginx-proxy` adında çalışan bir konteyneri durdurmak için `Stop` fonksiyonunun nasıl kullanılacağını gösteren bir örnek:

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerToStop := "my-nginx-proxy"

	log.Printf("'%s' konteyneri durduruluyor...", containerToStop)

	// pouch.Stop fonksiyonunu çağır
	err := pouch.Stop(containerToStop)
	if err != nil {
		log.Fatalf("Konteyner durdurulamadı: %v", err)
	}

	log.Printf("'%s' konteyneri başarıyla durduruldu!", containerToStop)
}
```