### `pouch.CopyToContainer` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `CopyToContainer` adında bir fonksiyon tanımlar. Bu fonksiyonun temel amacı, yerel makinedeki (host) bir dosyayı veya dizini, çalışan bir Docker konteynerinin içine kopyalamaktır. Bu işlemi, arka planda `docker cp` komutunu çalıştırarak gerçekleştirir.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// CopyToContainer, yerel makinedeki bir dosyayı veya dizini belirtilen Docker konteynerine kopyalar.
// Bu işlem, "docker cp" komutunu kullanarak yapılır.
func CopyToContainer(ID, HostPath, TargetPath string) error {
	// "docker cp [HostPath] [ContainerID]:[TargetPath]" komutunu oluştur.
	cmd := exec.Command("docker", "cp", HostPath, fmt.Sprintf("%s:%s", ID, TargetPath))
	
	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını içeren
		// daha açıklayıcı bir hata mesajı döndür.
		return fmt.Errorf("docker cp hatası: %v, çıktı: %s", err, string(out))
	}
	
	// İşlem başarılıysa nil (hata yok) döndür.
	return nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "cp", HostPath, fmt.Sprintf("%s:%s", ID, TargetPath))` satırı, Go'nun `os/exec` paketini kullanarak yeni bir komut nesnesi oluşturur. Bu komut, terminalde çalıştırılacak olan `docker cp [HostPath] [ContainerID]:[TargetPath]` komutuna denktir.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. Hem standart çıktıyı (stdout) hem de standart hatayı (stderr) tek bir byte dizisinde (`out`) yakalar. Eğer komut başarıyla çalışmazsa, bir `error` nesnesi (`err`) döndürür. Bu yöntem, hata ayıklama için çok yararlıdır, çünkü komut başarısız olduğunda Docker'ın verdiği hata mesajını da görebiliriz.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun çalışması sırasında herhangi bir hata olup olmadığını kontrol eder. Eğer bir hata varsa (`err` `nil` değilse), fonksiyon `fmt.Errorf` ile detaylı bir hata mesajı oluşturur. Bu mesaj, orijinal hatayı (`%v`, `err`) ve komutun ürettiği çıktıyı (`%s`, `string(out)`) içerir. Bu sayede, fonksiyonu çağıran kod, hatanın nedenini kolayca anlayabilir.

4.  **Başarılı Durum**:
    Eğer komut başarıyla tamamlanırsa, fonksiyon `nil` (hata yok) döndürerek işlemin başarılı olduğunu belirtir.

### Parametreler

*   `ID (string)`: Kopyalama yapılacak Docker konteynerinin kimliği (ID) veya adı.
*   `HostPath (string)`: Yerel makinede bulunan ve konteynere kopyalanacak olan dosyanın veya dizinin yolu.
*   `TargetPath (string)`: Dosyanın veya dizinin konteyner içinde kopyalanacağı hedef yol.

### Dönüş Değeri

*   `error`: Fonksiyon bir `error` değeri döndürür.
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   `docker cp` komutu başarısız olursa, komutun neden başarısız olduğuna dair bilgileri içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

Bu fonksiyonun doğru çalışabilmesi için, kodu çalıştıran sistemde **Docker'ın yüklü olması** ve `docker` komutunun sistemin `PATH` değişkeninde bulunması gerekmektedir.

### Kullanım Örneği

Yerel makinedeki `./config.yaml` dosyasını, `my-app-container` adlı konteynerin içindeki `/etc/app/` dizinine kopyalamak istediğimizi varsayalım.

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
	// Örneğin: "github.com/your-username/your-project/pouch"
)

func main() {
	containerID := "my-app-container"
	hostFile := "./config.yaml"
	containerPath := "/etc/app/"

	// pouch.CopyToContainer fonksiyonunu çağır
	err := pouch.CopyToContainer(containerID, hostFile, containerPath)
	if err != nil {
		log.Fatalf("Konteynere kopyalama başarısız oldu: %v", err)
	}

	log.Println("Dosya başarıyla konteynere kopyalandı!")
}
```