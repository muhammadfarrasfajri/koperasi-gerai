package utils

import (
	"regexp"
	"strings"
)

func ValidateNPWP(npwp string) (bool, string) {
    // 1. CEK KOSONG (Optionality)
    // Jika kosong, dianggap benar (karena optional)
    if npwp == "" {
        return true, "" 
    }
    cleanNPWP := strings.ReplaceAll(npwp, ".", "")
    cleanNPWP = strings.ReplaceAll(cleanNPWP, "-", "")

    // 3. CEK NUMERIK (Hanya boleh angka)
    isNumeric := regexp.MustCompile(`^[0-9]+$`).MatchString(cleanNPWP)
    if !isNumeric {
        return false, "NPWP must be in the form of numbers"
    }

    // 4. CEK PANJANG (15 atau 16 digit)
    length := len(cleanNPWP)
    if length != 15 && length != 16 {
        return false, "The length of the NPWP must be 15 or 16 digits"
    }

    // Lolos semua pengecekan
    return true, ""
}
