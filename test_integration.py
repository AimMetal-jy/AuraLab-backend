#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
集成测试脚本
测试WhisperX Flask服务和Go服务的集成
"""

import requests
import time
import os
#import json

def test_flask_service():
    """
    测试Flask服务
    """
    print("=== 测试Flask WhisperX服务 ===")
    
    # 测试健康检查
    try:
        response = requests.get('http://localhost:5000/health', timeout=5)
        if response.status_code == 200:
            print("✓ Flask服务健康检查通过")
            print(f"  响应: {response.json()}")
        else:
            print(f"✗ Flask服务健康检查失败: {response.status_code}")
            return False
    except requests.exceptions.RequestException as e:
        print(f"✗ Flask服务连接失败: {e}")
        return False
    
    return True

def test_go_service():
    """
    测试Go服务
    """
    print("\n=== 测试Go BlueLM服务 ===")
    
    # 测试TTS接口
    try:
        tts_data = {
            "mode": "human",
            "text": "这是一个测试",
            "vcn": "M24"
        }
        response = requests.post('http://localhost:8888/bluelm/tts', 
                               json=tts_data, timeout=10)
        if response.status_code == 200:
            print("✓ Go服务TTS接口可访问")
        else:
            print(f"✗ Go服务TTS接口失败: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"✗ Go服务连接失败: {e}")
        return False
    
    return True

def test_whisperx_integration():
    """
    测试WhisperX集成（如果有测试音频文件）
    """
    print("\n=== 测试WhisperX集成 ===")
    
    # 查找测试音频文件
    audio_file = None
    possible_paths = [
        r"C:\Users\Administrator\Downloads\Music\BBC_News.wav",
        "../BlueLM/audio_examples/English_Pod_30s.wav",
        "../audio_examples/English_Pod_30s.wav",
        "audio_examples/English_Pod_30s.wav"
    ]
    
    for path in possible_paths:
        if os.path.exists(path):
            audio_file = path
            break
    
    if not audio_file:
        print("⚠ 未找到测试音频文件，跳过集成测试")
        return True
    
    print(f"使用测试文件: {audio_file}")
    
    try:
        # 测试Go服务的whisperx接口
        with open(audio_file, 'rb') as f:
            files = {'file': f}
            response = requests.post('http://localhost:8888/whisperx', 
                                   files=files, timeout=300)
        
        if response.status_code == 200:
            result = response.json()
            print("✓ WhisperX集成测试成功")
            print(f"  处理结果: {result.get('message', 'No message')}")
            if 'data' in result:
                transcription = result['data'].get('transcription', [])
                print(f"  转录段数: {len(transcription)}")
        else:
            print(f"✗ WhisperX集成测试失败: {response.status_code}")
            print(f"  响应: {response.text}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"✗ WhisperX集成测试连接失败: {e}")
        return False
    
    return True

def main():
    """
    主测试函数
    """
    print("AuraLab后端服务集成测试")
    print("=" * 50)
    
    # 等待服务启动
    print("等待服务启动...")
    time.sleep(2)
    
    # 执行测试
    flask_ok = test_flask_service()
    go_ok = test_go_service()
    
    if flask_ok and go_ok:
        integration_ok = test_whisperx_integration()
    else:
        print("\n⚠ 基础服务测试失败，跳过集成测试")
        integration_ok = False
    
    # 输出测试结果
    print("\n" + "=" * 50)
    print("测试结果汇总:")
    print(f"Flask服务: {'✓ 通过' if flask_ok else '✗ 失败'}")
    print(f"Go服务: {'✓ 通过' if go_ok else '✗ 失败'}")
    print(f"集成测试: {'✓ 通过' if integration_ok else '✗ 失败'}")
    
    if flask_ok and go_ok:
        print("\n🎉 所有基础服务正常运行！")
        print("服务地址:")
        print("  - WhisperX Flask: http://localhost:5000")
        print("  - BlueLM Go: http://localhost:8888")
        print("  - 统一接口: http://localhost:8888/whisperx")
    else:
        print("\n❌ 部分服务未正常运行，请检查服务启动状态")

if __name__ == '__main__':
    main()