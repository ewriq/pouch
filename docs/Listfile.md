
### `pouch.ListFiles` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `ListFiles` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, çalışan bir Docker konteynerinin içindeki belirli bir dizinin içeriğini (dosya ve klasörleri) listelemektir. Bu işlemi, `docker exec` komutu aracılığıyla konteyner içinde `ls -l` komutunu çalıştırarak yapar ve çıktının her bir satırını bir string dizisi elemanı olarak döndürür.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// ListFiles, belirtilen bir Docker konteynerindeki bir dizinin içeriğini
// 'ls -l' formatında listeler.
// Her bir satır, bir string dizisinin elemanı olarak döndürülür.
func ListFiles(containerID, containerPath string) ([]string, error) {
	// Komutu oluştur: docker exec [ID] ls -l [dizin-yolu]
	// "ls -l" komutu detaylı (long) bir liste formatı sağlar.
	cmd := exec.Command("docker", "exec", containerID, "ls", "-l", containerPath)
	
	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return nil, fmt.Errorf("Dosya listeleme hatası: %v, %s", err, string(out))
	}

	// Komutun çıktısındaki baştaki ve sondaki boşlukları temizle,
	// ardından çıktıyı satır satır bölerek bir string dizisi oluştur.
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	// Elde edilen satır dizisini döndür.
	return lines, nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "exec", containerID, "ls", "-l", containerPath)` satırı, çalıştırılacak komutu oluşturur. Bu komut, belirtilen konteynerin (`containerID`) içinde `ls -l [containerPath]` komutunu çalıştırır. `-l` bayrağı, dosyaların izinleri, sahibi, boyutu, son değiştirilme tarihi gibi detaylı bilgileri içeren uzun bir liste formatı sağlar.

2.  **Komutun Çalıştırılması ve Hata Yönetimi**:
    `cmd.CombinedOutput()` fonksiyonu, komutu çalıştırır ve hem standart çıktıyı (dosya listesi) hem de standart hatayı (hata mesajları) tek bir `[]byte` dizisinde toplar. Eğer `docker exec` bir hata verirse (örneğin konteyner çalışmıyorsa veya belirtilen dizin yoksa), `err` değişkeni `nil` olmaz ve fonksiyon, detaylı bir hata mesajıyla birlikte `nil` bir dizi döndürür.

3.  **Çıktının İşlenmesi**:
    Bu fonksiyonun en önemli adımı burasıdır.
    *   `string(out)`: Ham byte çıktısı bir metne (string) dönüştürülür.
    *   `strings.TrimSpace(...)`: Bu metnin başındaki ve sonundaki gereksiz boşluklar ve yeni satır karakterleri temizlenir.
    *   `strings.Split(..., "\n")`: Temizlenmiş metin, yeni satır karakterlerinden (`\n`) bölünerek bir string dizisi (`[]string`) haline getirilir. Sonuç olarak, `ls -l` komutunun her bir satırı, bu dizinin ayrı bir elemanı olur.

4.  **Başarılı Dönüş**:
    İşlem başarılı olursa, dosya listesinin her satırını içeren `lines` dizisi ve `nil` bir hata değeri döndürülür.

### Parametreler

*   `containerID (string)`: İçinde dosya listelenecek olan Docker konteynerinin kimliği (ID) veya adı.
*   `containerPath (string)`: Konteynerin **içinde** listelenecek olan dizinin yolu (örneğin, `/app`, `/var/log` vb.).

### Dönüş Değeri

*   `[]string`: İşlem başarılı olursa, `ls -l` komutunun çıktısındaki her bir satırı içeren bir string dizisi döner. Dizinin ilk satırı genellikle "total <blok_sayısı>" şeklinde bir bilgi içerir. Her satır, dosya izinleri, sahibi, boyutu gibi bilgileri barındırır.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner çalışmıyorsa, dizin mevcut değilse veya başka bir `docker` hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Önemli Notlar

*   Fonksiyonun çalışması için sistemde **Docker'ın yüklü** ve çalışır durumda olması gerekmektedir.
*   Hedef `containerID` ile belirtilen konteynerin **çalışıyor olması** zorunludur.

### Kullanım Örneği

`nginx-server` adlı bir konteynerin `/etc/nginx/conf.d` dizinindeki yapılandırma dosyalarını listelemek için:

```go
package main

import (
	"fmt"
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerID := "nginx-server"
	directoryPath := "/etc/nginx/conf.d"

	// pouch.ListFiles fonksiyonunu çağırarak dizin içeriğini al
	fileList, err := pouch.ListFiles(containerID, directoryPath)
	if err != nil {
		log.Fatalf("Konteynerdeki dosyalar listelenemedi: %v", err)
	}

	fmt.Printf("'%s' konteynerindeki '%s' dizininin içeriği:\n", containerID, directoryPath)
	
	// Dönen dizi 'ls -l' çıktısının ham halidir. İlk satır genellikle 'total' bilgisini içerir.
	// İsterseniz bu satırı atlayabilirsiniz.
	for i, line := range fileList {
		if i == 0 && strings.HasPrefix(line, "total") {
			fmt.Printf("(Toplam blok bilgisi: %s)\n", line)
			continue
		}
		fmt.Println(line)
	}
}
```