
### `pouch.Exec` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Exec` adında genel amaçlı bir fonksiyon tanımlar. Fonksiyonun temel amacı, çalışan bir Docker konteynerinin içinde istenilen herhangi bir komutu çalıştırmak ve bu komutun ürettiği çıktıyı (hem standart çıktı hem de standart hata) bir metin olarak geri döndürmektir. Bu, `docker exec` komutunu programatik olarak kullanmak için güçlü ve esnek bir yöntem sunar.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// Exec, belirtilen ID'ye sahip bir Docker konteyneri içinde verilen komutu çalıştırır.
// Komutun standart çıktı (stdout) ve standart hatasını (stderr) birleştirilmiş
// bir string olarak döndürür.
func Exec(id string, command []string) (string, error) {
	// "docker exec [id]" komutunun temel argümanlarını oluştur ve 
	// çalıştırılacak komutu bu argümanların sonuna ekle.
	args := append([]string{"exec", id}, command...)
	
	// Docker komutunu oluştur.
	cmd := exec.Command("docker", args...)

	// Komutu çalıştır ve hem standart çıktıyı hem de standart hatayı yakala.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem Go'nun hata nesnesini hem de komutun çıktısını
		// içeren detaylı bir hata mesajı döndür.
		return "", fmt.Errorf("exec hatası: %v\nçıktı: %s", err, out)
	}

	// Başarılı olursa, komutun çıktısını string olarak döndür.
	return string(out), nil
}
```

### Fonksiyon Detayları

1.  **Argümanların Birleştirilmesi**:
    `args := append([]string{"exec", id}, command...)` satırı, fonksiyonun temelini oluşturur.
    *   İlk olarak `[]string{"exec", id}` dilimi ile `docker exec [containerID]` komutunun başlangıcı hazırlanır.
    *   Daha sonra `append` fonksiyonu ve `...` operatörü kullanılarak, `command` parametresi olarak gelen `string` diliminin tüm elemanları (`çalıştırılacak_komut` ve `argümanları`) bu listeye eklenir. Örneğin, `command` `[]string{"ls", "-la", "/app"}` ise, `args` dilimi `[]string{"exec", "my-container", "ls", "-la", "/app"}` haline gelir. Bu yapı, komutları güvenli bir şekilde (shell injection riski olmadan) oluşturmayı sağlar.

2.  **Komutun Çalıştırılması**:
    `cmd := exec.Command("docker", args...)` ile çalıştırılacak tam komut nesnesi yaratılır. `cmd.CombinedOutput()` fonksiyonu bu komutu çalıştırır ve en önemli özelliklerinden biri olan standart çıktı (stdout) ile standart hatayı (stderr) tek bir byte dizisinde (`out`) birleştirir. Bu, komut başarısız olduğunda bile hata mesajlarını yakalamayı garanti eder.

3.  **Hata Yönetimi**:
    `if err != nil` bloğu, komutun sıfırdan farklı bir çıkış koduyla sonlanması gibi hata durumlarını yakalar. Bu durumda, `fmt.Errorf` ile oldukça bilgilendirici bir hata mesajı oluşturulur. Bu mesaj hem Go seviyesindeki hatayı (`err`) hem de komutun kendisinin ürettiği çıktıyı (`out`) içerir. Bu sayede, hata ayıklama süreci büyük ölçüde kolaylaşır.

4.  **Başarılı Çıktı**:
    Komut başarıyla çalışırsa (çıkış kodu 0 ise), yakalanan çıktı (`out` byte dizisi) `string(out)` ifadesiyle bir metne dönüştürülür ve `nil` hatasıyla birlikte döndürülür.

### Parametreler

*   `id (string)`: İçinde komut çalıştırılacak olan Docker konteynerinin kimliği (ID) veya adı.
*   `command ([]string)`: Konteyner içinde çalıştırılacak olan komut ve argümanlarını içeren bir `string` dilimi. Dilimin ilk elemanı komutun kendisi, sonraki elemanlar ise argümanları olmalıdır.

### Dönüş Değeri

*   `string`: Komutun çalışması sonucu ortaya çıkan birleştirilmiş standart çıktı ve standart hata metni.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   `docker exec` komutu başarısız olursa (örneğin konteyner çalışmıyorsa veya komut hata verirse), detaylı bilgi içeren bir hata nesnesi döndürülür.

### Önemli Notlar

*   Bu fonksiyonun çalışabilmesi için sistemde **Docker'ın kurulu** ve `docker` komutunun `PATH` içinde erişilebilir olması gerekir.
*   Hedef konteynerin `id` ile belirtilen kimlikle **çalışıyor durumda** olması zorunludur.

### Kullanım Örneği

`my-app` adlı bir konteynerin içindeki tüm ortam değişkenlerini (`env` komutu ile) listelemek için:

```go
package main

import (
	"fmt"
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerID := "my-app"
	commandToRun := []string{"env"} // Sadece komut, argüman yok.

	// Konteynerdeki /etc dizinini listelemek için:
	// commandToRun := []string{"ls", "-l", "/etc"}

	output, err := pouch.Exec(containerID, commandToRun)
	if err != nil {
		log.Fatalf("Konteyner içinde komut çalıştırılamadı: %v", err)
	}

	fmt.Printf("'%s' konteynerinden gelen çıktı:\n%s", containerID, output)
}
```