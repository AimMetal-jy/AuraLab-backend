package utils

import (
	"encoding/binary"
	"fmt"
	"os"
)

func PcmtoWav(pcmData []byte, filename string, channels, bitsPerSample, sampleRate int) error {
	if bitsPerSample%8 != 0 {
		return fmt.Errorf("bits %% 8 must == 0. now bits: %d", bitsPerSample)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 计算参数
	sampleWidth := bitsPerSample / 8
	dataSize := len(pcmData)
	byteRate := sampleRate * channels * sampleWidth
	blockAlign := channels * sampleWidth
	fileSize := 36 + dataSize

	// RIFF头 (12字节)
	file.Write([]byte("RIFF"))
	binary.Write(file, binary.LittleEndian, uint32(fileSize))
	file.Write([]byte("WAVE"))

	// fmt子块 (24字节)
	file.Write([]byte("fmt "))
	binary.Write(file, binary.LittleEndian, uint32(16))                // fmt子块大小
	binary.Write(file, binary.LittleEndian, uint16(1))                 // 音频格式(PCM)
	binary.Write(file, binary.LittleEndian, uint16(channels))           // 声道数
	binary.Write(file, binary.LittleEndian, uint32(sampleRate))         // 采样率
	binary.Write(file, binary.LittleEndian, uint32(byteRate))           // 字节率
	binary.Write(file, binary.LittleEndian, uint16(blockAlign))         // 块对齐
	binary.Write(file, binary.LittleEndian, uint16(bitsPerSample))      // 位深度

	// data子块 (8字节头 + 数据)
	file.Write([]byte("data"))
	binary.Write(file, binary.LittleEndian, uint32(dataSize))
	file.Write(pcmData)

	return nil
}