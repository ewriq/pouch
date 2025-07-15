
### `pouch.Logs` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Logs` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, belirtilen bir Docker konteynerinin başlangıcından itibaren biriktirdiği tüm log kayıtlarını (standart çıktı ve standart hata) almaktır. Bu işlemi, arka planda `docker logs` komutunu çalıştırarak gerçekleştirir ve komutun çıktısını ham metin olarak döndürür.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// Logs, belirtilen ID'ye sahip Docker konteynerinin tüm log kayıtlarını getirir.
func Logs(id string) (string, error) {
	// "docker logs [id]" komutunu oluştur.
	cmd := exec.Command("docker", "logs", id)

	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	// docker logs komutu zaten bu iki akışı birleştirerek verir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return "", fmt.Errorf("logs hatası: %v, çıktı: %s", err, string(out))
	}

	// İşlem başarılıysa, logları string olarak döndür.
	return string(out), nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "logs", id)` satırı, terminalde `docker logs [containerID]` komutunu çalıştırmak için bir komut nesnesi oluşturur. Bu komut, bir konteynerin çalıştığı süre boyunca standart çıktı (`stdout`) ve standart hata (`stderr`) akışlarına yazdığı tüm verileri getirir.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. `docker logs` komutu doğası gereği hem `stdout` hem de `stderr` loglarını bir arada sunduğu için, `CombinedOutput` bu çıktıyı eksiksiz bir şekilde yakalamak için ideal bir yöntemdir. Komutun tüm çıktısı `out` adlı bir byte dizisinde toplanır.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun çalışması sırasında bir sorun olup olmadığını denetler. En yaygın hatalar arasında, belirtilen ID'ye sahip bir konteynerin bulunamaması veya Docker servisinin (daemon) çalışmıyor olması yer alır. Bir hata oluştuğunda, fonksiyon hem Go seviyesindeki hatayı (`err`) hem de Docker CLI'nin ürettiği hata mesajını (`out`) içeren zengin bir hata mesajıyla geri döner.

4.  **Başarılı Durum**:
    Komut başarıyla çalışırsa, `string(out)` ifadesiyle byte dizisi halindeki log verisi bir metne (string) dönüştürülür ve `nil` hata değeri ile birlikte fonksiyonun çağrıldığı yere döndürülür.

### Parametreler

*   `id (string)`: Logları alınacak olan Docker konteynerinin kimliği (ID) veya adı.

### Dönüş Değeri

*   `string`: İşlem başarılı olursa, konteynerin tüm birikmiş loglarını içeren ham metin.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner bulunamazsa veya Docker servisiyle ilgili bir sorun yaşanırsa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

Bu fonksiyonun çalışabilmesi için, kodu yürüten sistemde **Docker'ın yüklü** olması ve `docker` komutunun sistemin `PATH` değişkeni üzerinden erişilebilir olması gerekmektedir.

### Kullanım Örneği

`my-api-server` adlı bir konteynerin tüm loglarını alıp konsola yazdırmak için `Logs` fonksiyonunun nasıl kullanılacağını gösteren bir örnek:

```go
package main

import (
	"fmt"
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerID := "my-api-server"

	// pouch.Logs fonksiyonunu çağırarak konteynerin loglarını al
	logOutput, err := pouch.Logs(containerID)
	if err != nil {
		log.Fatalf("Konteyner logları alınamadı: %v", err)
	}

	fmt.Printf("'%s' Konteynerinin Logları:\n", containerID)
	fmt.Println("---------------------------------")
	fmt.Println(logOutput)
	fmt.Println("---------------------------------")
}
```