### `pouch.List` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `List` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, sistemde bulunan **tüm** Docker konteynerlerini (hem çalışanları hem de durdurulmuş olanları) listelemektir. Bu işlem, arka planda `docker ps -a` komutunu çalıştırarak gerçekleştirilir ve komutun ham metin çıktısını bir string olarak döndürür.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// List, sistemdeki tüm Docker konteynerlerini (çalışan ve durdurulmuş) listeler.
// Sonuç olarak "docker ps -a" komutunun ham çıktısını döndürür.
func List() (string, error) {
	// "docker ps -a" komutunu oluştur. `-a` bayrağı durdurulmuş olanlar dahil tüm
	// konteynerleri listelemeyi sağlar.
	cmd := exec.Command("docker", "ps", "-a") 
	
	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return "", fmt.Errorf("listeleme hatası: %v, çıktı: %s", err, string(out))
	}
	
	// İşlem başarılıysa, komutun çıktısını string olarak döndür.
	return string(out), nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "ps", "-a")` satırı, terminalde `docker ps -a` komutunu çalıştırmak için bir komut nesnesi hazırlar.
    *   `docker ps`: Genellikle çalışan konteynerleri listeler.
    *   `-a` (`--all`): Bu bayrak, komutun durdurulmuş konteynerleri de listeye dahil etmesini sağlar. Bu fonksiyon, bu bayrağı standart olarak kullanarak sistemdeki tüm konteynerler hakkında bilgi verir.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. Komutun standart çıktısını (genellikle konteyner tablosu) ve standart hatasını (olası hata mesajları) tek bir byte dizisinde (`out`) birleştirir. Bu, özellikle Docker servisi çalışmıyorsa veya başka bir sorun varsa, hata mesajlarını yakalamak için çok kullanışlıdır.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun başarıyla çalışıp çalışmadığını kontrol eder. Örneğin, Docker daemon (servisi) çalışmıyorsa, `exec` paketi bir hata döndürür. Bu durumda fonksiyon, hatanın nedenini anlamayı kolaylaştırmak için hem orijinal Go hatasını (`err`) hem de Docker'dan gelen çıktıyı (`out`) içeren detaylı bir hata mesajı oluşturur.

4.  **Başarılı Durum**:
    Komut hatasız bir şekilde tamamlandığında, `string(out)` ifadesiyle byte dizisi halindeki ham çıktı bir metne (string) dönüştürülür ve `nil` hata değeri ile birlikte geri döndürülür. Döndürülen metin, terminalde `docker ps -a` çalıştırdığınızda göreceğinizle birebir aynıdır.

### Parametreler

*   Bu fonksiyon herhangi bir parametre **almaz**.

### Dönüş Değeri

*   `string`: İşlem başarılı olursa, tüm konteynerleri listeleyen ve terminalde görüneceği formatta olan ham metin çıktısı.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer Docker servisine ulaşılamazsa veya başka bir `docker` hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

Bu fonksiyonun beklendiği gibi çalışması için, Go kodunun çalıştığı ortamda **Docker'ın yüklü** olması ve `docker` komutunun sistemin `PATH` değişkeninde bulunması zorunludur.

### Kullanım Örneği

Aşağıdaki örnek, sistemdeki tüm konteynerleri listelemek ve sonucu konsola yazdırmak için `List` fonksiyonunun nasıl kullanılacağını gösterir.

```go
package main

import (
	"fmt"
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	// pouch.List fonksiyonunu çağırarak tüm konteynerlerin listesini al
	containerList, err := pouch.List()
	if err != nil {
		log.Fatalf("Konteynerler listelenemedi: %v", err)
	}

	fmt.Println("Sistemdeki Tüm Docker Konteynerleri:")
	// Dönen ham metni doğrudan ekrana yazdır
	fmt.Println(containerList)
}
```