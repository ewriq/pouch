### `pouch.ContainerStats` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `ContainerStats` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, **çalışan** bir Docker konteynerinin anlık kaynak kullanım istatistiklerini (CPU, bellek, ağ ve disk G/Ç) verimli bir şekilde almaktır. Bu işlemi, `docker stats` komutunu özel bir formatlama seçeneğiyle çalıştırarak ve dönen çıktıyı programatik olarak işleyerek bir `map` veri yapısı halinde döndürür.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// ContainerStats, çalışan bir konteynerin anlık kaynak kullanımını
// (CPU, Bellek, Ağ, Disk) bir map olarak döndürür.
func ContainerStats(containerID string) (map[string]string, error) {
	// Dönecek olan map'i başlat.
	stats := make(map[string]string)

	// "docker stats" komutunu özel formatlama seçenekleriyle oluştur.
	// --no-stream: Komutun sürekli veri akışı yapmak yerine anlık bir görüntü alıp çıkmasını sağlar.
	// --format: Çıktıyı bizim istediğimiz gibi "KEY=VALUE, KEY=VALUE" formatında yapılandırır.
	cmd := exec.Command("docker", "stats", "--no-stream", "--format",
		"CPU={{.CPUPerc}}, MEM={{.MemUsage}}, NET={{.NetIO}}, DISK={{.BlockIO}}", containerID)

	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("stats hatası: %v\nçıktı: %s", err, out)
	}

	// Komutun çıktısını ayrıştırma (parsing) işlemi.
	// Örnek çıktı: "CPU=0.15%, MEM=12.5MiB / 1.95GiB, NET=648B / 0B, DISK=0B / 8.19kB"
	// 1. Çıktıyı ", " karakterine göre bölerek her bir istatistik parçasını ayır.
	lines := strings.Split(strings.TrimSpace(string(out)), ", ")
	for _, line := range lines {
		// 2. Her parçayı "=" karakterine göre ikiye ayır (anahtar ve değer).
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			// 3. Anahtarı temizle ve küçük harfe çevir, değeri temizle.
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			val := strings.TrimSpace(parts[1])
			// 4. İstatistikleri map'e ekle.
			stats[key] = val
		}
	}

	return stats, nil
}
```

### Fonksiyon Detayları

Bu fonksiyonun çalışması iki ana aşamadan oluşur: veri toplama ve veri işleme.

1.  **Veri Toplama (Komutun Akıllıca Oluşturulması)**:
    Fonksiyon, basit bir `docker stats` komutu çalıştırmak yerine, onu belirli bayraklarla zenginleştirir:
    *   `--no-stream`: `docker stats` komutu normalde sürekli bir veri akışı sağlar. Bu bayrak, komutun sadece **anlık** bir kaynak kullanım verisi alıp hemen sonlanmasını sağlar. Bu, programatik kullanım için kritiktir.
    *   `--format`: Bu bayrak, çıktının formatını Go'nun şablonlama (templating) dilini kullanarak özelleştirmemize olanak tanır. Burada, çıktı kasıtlı olarak `KEY={{.Value}}` şeklinde ve virgülle ayrılmış bir yapıya dönüştürülür. Bu sayede, sonradan işlenmesi çok kolay olan, tahmin edilebilir bir metin elde edilir.

2.  **Veri İşleme (Çıktının Ayrıştırılması - Parsing)**:
    `docker` komutundan dönen `CPU=..., MEM=..., NET=...` şeklindeki ham metin, Go'nun `strings` paketi kullanılarak işlenir:
    *   `strings.TrimSpace` ile çıktının başındaki ve sonundaki olası boşluklar temizlenir.
    *   `strings.Split(..., ", ")` ile metin, virgül ve boşluk karakterlerinden bölünerek `["CPU=...", "MEM=...", ...]` şeklinde bir diziye dönüştürülür.
    *   Bir `for` döngüsü içinde, her bir dizi elemanı bu sefer `=` karakterinden `strings.SplitN` ile ikiye bölünerek anahtar (key) ve değer (value) elde edilir. `SplitN` kullanmak, değer kısmında `=` karakteri olması ihtimaline karşı bir güvenlik önlemidir.
    *   Anahtar (`key`) küçük harfe çevrilir ve hem anahtar hem de değer (`val`) `TrimSpace` ile temizlenir.
    *   Son olarak, temizlenmiş anahtar ve değer çifti `stats` adlı `map`'e eklenir.

### Parametreler

*   `containerID (string)`: İstatistikleri alınacak olan, **çalışır durumdaki** Docker konteynerinin kimliği (ID) veya adı.

### Dönüş Değeri

*   `map[string]string`: İşlem başarılı olursa, aşağıdaki anahtarları ve karşılık gelen değerleri içeren bir `map` döner:
    *   `"cpu"`: CPU kullanım yüzdesi.
    *   `"mem"`: Bellek kullanımı (ör: "12.5MiB / 1.95GiB").
    *   `"net"`: Ağ G/Ç (giriş/çıkış) (ör: "648B / 0B").
    *   `"disk"`: Blok (disk) G/Ç (ör: "0B / 8.19kB").
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner bulunamazsa, çalışmıyorsa veya başka bir Docker hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar ve Önemli Notlar

*   Fonksiyonun çalışması için sistemde **Docker'ın yüklü** olması gerekmektedir.
*   `docker stats` komutu sadece **çalışan konteynerler** için veri döndürdüğünden, hedef konteynerin çalışır durumda olması zorunludur.

### Kullanım Örneği

`my-service` adındaki bir konteynerin anlık kaynak kullanımını alıp yazdırma:

```go
package main

import (
	"fmt"
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerID := "my-service"

	stats, err := pouch.ContainerStats(containerID)
	if err != nil {
		log.Fatalf("Konteyner istatistikleri alınamadı: %v", err)
	}

	fmt.Printf("'%s' Konteynerinin Anlık İstatistikleri:\n", containerID)
	fmt.Printf("  - CPU Kullanımı: %s\n", stats["cpu"])
	fmt.Printf("  - Bellek Kullanımı: %s\n", stats["mem"])
	fmt.Printf("  - Ağ G/Ç: %s\n", stats["net"])
	fmt.Printf("  - Disk G/Ç: %s\n", stats["disk"])
}
```