import BasePageController from "../base/BasePageController.js";

class MainPageController extends BasePageController {
    constructor(router, pageHolder) {
        super(router, pageHolder);
    }

    pageInit() {
        super.pageInit();
        
        // 차트 설정 (전역 변수로 선언하여 쉽게 수정 가능)
        window.chartConfig = {
            // 차트 선 색상
            lineColor: '#1b02fd',           // 메인 차트 선 색상
            gridColor: '#333',               // 그리드 색상
            currentPriceColor: '#00c853',    // 현재 가격 점 색상
            currentPriceLineColor: '#ffffff', // 현재가 수평선 색상 (하얀색)
            currentPriceTextColor: '#ffffff', // 현재가 텍스트 색상 (하얀색)
            expirationLineColor: '#ff9800',  // 만기 시간 선 색상
            upperRangeLineColor: '#ff1744',  // 상단 범위 선 색상
            lowerRangeLineColor: '#00e676',   // 하단 범위 선 색상
            upColor: '#00ff00',               // 상승 캔들 색상 (초록)
            downColor: '#ff0000',             // 하락 캔들 색상 (빨강)
            candleWidth: 0.8,                 // 캔들 너비 비율 (0~1)
            
            // Y축 비율 설정
            yAxisRangePercent: 0.007,         // Y축 범위 (현재 가격 기준 ±2%)
            
            // X축 비율 설정
            pastDataPoints: 100,             // 과거 데이터 포인트 수
            futureDataPoints: 60,             // 미래 데이터 포인트 수 (5초 단위)

            timeIntervalSeconds: 5            // 시간 간격 (초)
        };
        
        this.initTradeTypeSelector();
        this.createChart();
    }

    initTradeTypeSelector() {
        const selector = document.getElementById('trade-type-selector');
        const display = document.getElementById('trade-type-display');
        const dropdown = document.getElementById('trade-type-dropdown');
        
        if (!selector || !display || !dropdown) {
            console.error('요소를 찾을 수 없습니다');
            return;
        }
        
        const options = dropdown.querySelectorAll('.trade-option');
        
        if (options.length === 0) {
            console.error('옵션을 찾을 수 없습니다');
            return;
        }

        let currentType = null;
        let closeTimeout = null;

        // 클릭으로 드롭다운 열기/닫기
        selector.addEventListener('click', (e) => {
            e.stopPropagation();
            if (dropdown.classList.contains('hidden')) {
                dropdown.classList.remove('hidden');
            } else {
                dropdown.classList.add('hidden');
            }
        });
        
        // 외부 클릭 시 닫기
        document.addEventListener('click', (e) => {
            if (!selector.contains(e.target)) {
                dropdown.classList.add('hidden');
            }
        });

        options.forEach(option => {
            option.addEventListener('mousedown', (e) => {
                e.preventDefault();
                e.stopPropagation();
            });
            
            option.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                const type = option.dataset.type;
                currentType = type;
                
                options.forEach(opt => opt.classList.remove('selected'));
                option.classList.add('selected');
                
                if (type === 'updown') {
                    display.textContent = 'UP/DOWN';
                    window.currentTradeType = 'updown';
                    // Touch/NoTouch 범위 선 숨기기
                    if (window.touchRangeLines) {
                        window.touchRangeLines.upper = null;
                        window.touchRangeLines.lower = null;
                    }
                } else if (type === 'touchnotouch') {
                    display.textContent = 'Touch/NoTouch';
                    window.currentTradeType = 'touchnotouch';
                    
                    // touchRangeLines가 없으면 초기화
                    if (!window.touchRangeLines) {
                        window.touchRangeLines = {
                            upper: null,
                            lower: null,
                            isDragging: false,
                            dragLine: null,
                            dragStartY: 0,
                            dragStartUpper: 0,
                            dragStartLower: 0
                        };
                    }
                    
                    // Touch/NoTouch 범위 선 초기화 (현재 데이터 범위 기준)
                    if (window.data && window.data.length > 0) {
                        const min = Math.min(...window.data);
                        const max = Math.max(...window.data);
                        const dataRange = max - min;
                        const center = (min + max) / 2;
                        
                        // 데이터 범위의 약 20%를 범위 선 간격으로 설정
                        const lineRange = dataRange * 0.2;
                        window.touchRangeLines.upper = center + lineRange;
                        window.touchRangeLines.lower = center - lineRange;
                        
                    } else if (window.currentPrice > 0) {
                        // 데이터가 없으면 현재 가격 기준 ±2%
                            
                        const range = window.currentPrice * 0.02;
                        window.touchRangeLines.upper = window.currentPrice + range;
                        window.touchRangeLines.lower = window.currentPrice - range;
                    }
                }
                
                // 차트 다시 그리기
                if (window.drawChart) {
                    window.drawChart();
                }
                
                dropdown.classList.add('hidden');
            });
        });
    }

    createChart() {
        
        const chartContainer = document.getElementById('chart');
        if (!chartContainer) {
            console.error('차트 컨테이너를 찾을 수 없습니다.');
            return;
        }

        // Canvas 생성
        const canvas = document.createElement('canvas');
        canvas.width = chartContainer.clientWidth;
        canvas.height = chartContainer.clientHeight || 400;
        canvas.style.width = '100%';
        canvas.style.height = '100%';
        chartContainer.appendChild(canvas);
        
        const ctx = canvas.getContext('2d');
        
        // 실시간 데이터 배열 (전역으로 노출)
        // 캔들스틱 데이터: {open, high, low, close, time}
        const data = [];
        window.data = data;
        const maxDataPoints = window.chartConfig ? window.chartConfig.pastDataPoints : 100;
        let currentPrice = 0;
        window.currentPrice = currentPrice;
        
        // 미래 데이터 배열 (예측용, 5초 단위)
        const futureData = [];
        const futureDataPoints = window.chartConfig ? window.chartConfig.futureDataPoints : 60;
        
        // 만기 시간 설정 (데이터 인덱스 기준, 기본값: 현재 위치에서 60개 뒤)
        let expirationIndex = 60;
        let isDragging = false;
        let startPrice = 0; // 거래 시작 가격
        
        // Touch/NoTouch 범위 선 (Y축 가격 기준) - 전역으로 노출
        if (!window.touchRangeLines) {
            window.touchRangeLines = {
                upper: null, // 상단 선 가격
                lower: null,  // 하단 선 가격
                isDragging: false,
                dragLine: null, // 'upper' 또는 'lower'
                dragStartY: 0,
                dragStartUpper: 0,
                dragStartLower: 0,
                dragStartCenter: 0, // 드래그 시작 시 중심점
                dragStartRange: 0  // 드래그 시작 시 범위
            };
        }
        
        // 현재 선택된 거래 타입 (전역으로 노출)
        if (window.currentTradeType === undefined) {
            window.currentTradeType = null;
        }
        
        // Binance API에서 과거 데이터 가져오기 (초기 차트 채우기)
        const fetchHistoricalData = async () => {
            try {
                // 현재 가격 가져오기
                const priceResponse = await fetch('https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT');
                const priceResult = await priceResponse.json();
                const price = parseFloat(priceResult.price);
                
                currentPrice = price;
                window.currentPrice = currentPrice;
                startPrice = currentPrice;
                
                // 1분 캔들 데이터 가져오기 (OHLC)
                const klinesResponse = await fetch('https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1m&limit=100');
                const klines = await klinesResponse.json();
                
                // 캔들스틱 데이터로 변환
                klines.forEach(kline => {
                    data.push({
                        open: parseFloat(kline[1]),   // 시가
                        high: parseFloat(kline[2]),   // 고가
                        low: parseFloat(kline[3]),    // 저가
                        close: parseFloat(kline[4]),  // 종가
                        time: kline[0]                // 시간
                    });
                });
                
                // 만기 시간을 현재 위치에서 60개 뒤로 설정 (5분 = 60 * 5초)
                const currentIndex = maxDataPoints;
                expirationIndex = currentIndex + 60;
                
                // Touch/NoTouch가 선택되어 있고 범위 선이 없으면 초기화 (현재 금액에 맞춤)
                if (window.currentTradeType === 'touchnotouch' && window.touchRangeLines && 
                    (window.touchRangeLines.upper === null || window.touchRangeLines.lower === null)) {
                    const lastCandle = data.length > 0 ? data[data.length - 1] : null;
                    const price = currentPrice > 0 ? currentPrice : (lastCandle ? (lastCandle.close || lastCandle) : 87000);
                    const range = price * 0.05;
                    window.touchRangeLines.upper = price + range;
                    window.touchRangeLines.lower = price - range;
                }
                
                drawChart();
            } catch (error) {
                console.error('과거 데이터 가져오기 실패:', error);
                // 샘플 데이터로 채우기
                const basePrice = 87000;
        for (let i = 0; i < 50; i++) {
                    data.push(basePrice + (Math.random() - 0.5) * 1000);
                }
                currentPrice = data[data.length - 1];
                drawChart();
            }
        };
        
        // Binance API에서 현재 가격 가져오기
        const fetchPrice = async () => {
            try {
                // BTC/USDT 가격 가져오기
                const response = await fetch('https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT');
                const result = await response.json();
                const price = parseFloat(result.price);
                
                currentPrice = price;
                window.currentPrice = currentPrice;
                
                // 마지막 캔들의 종가를 현재 가격으로 업데이트
                if (data.length > 0) {
                    const lastCandle = data[data.length - 1];
                    if (lastCandle.close !== undefined) {
                        // 캔들 객체인 경우
                        lastCandle.close = price;
                        lastCandle.high = Math.max(lastCandle.high, price);
                        lastCandle.low = Math.min(lastCandle.low, price);
                    } else {
                        // 숫자인 경우 (기존 데이터 호환)
                        data[data.length - 1] = price;
                    }
                } else {
                    // 데이터가 없으면 새 캔들 생성
                    data.push({
                        open: price,
                        high: price,
                        low: price,
                        close: price,
                        time: Date.now()
                    });
                }
                
                // Touch/NoTouch가 선택되어 있고 범위 선이 없으면 초기화 (현재 금액에 맞춤)
                if (window.currentTradeType === 'touchnotouch' && window.touchRangeLines && 
                    (window.touchRangeLines.upper === null || window.touchRangeLines.lower === null)) {
                    const lastCandle = data.length > 0 ? data[data.length - 1] : null;
                    const price = currentPrice > 0 ? currentPrice : (lastCandle ? (lastCandle.close || lastCandle) : 87000);
                    const range = price * 0.05;
                    window.touchRangeLines.upper = price + range;
                    window.touchRangeLines.lower = price - range;
                }
                
                // 가격 표시 업데이트
                const strikeDiv = document.querySelector('#trade-panel > div:nth-child(3)');
                if (strikeDiv) {
                    const priceDiv = strikeDiv.querySelector('div:last-child');
                    if (priceDiv && priceDiv.textContent.includes('Price')) {
                        priceDiv.innerHTML = `Price: <span style="color: #00c853; font-weight: 700;">${price.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>`;
                    }
                }
                
                drawChart();
            } catch (error) {
                console.error('가격 가져오기 실패:', error);
                // 에러 발생 시 마지막 가격 유지
                if (data.length > 0) {
                    data.push(data[data.length - 1]);
                    if (data.length > maxDataPoints) {
                        data.shift();
                    }
                    drawChart();
                }
            }
        };
        
        const drawChart = () => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            
            // 배경
            ctx.fillStyle = '#000000';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            
            // Y축 레이블 영역 (우측)
            const yAxisWidth = 80;
            const xAxisHeight = 50;
            
            // 차트 영역
            const chartX = 0;
            const chartY = 0;
            const chartWidth = canvas.width - yAxisWidth;
            const chartHeight = canvas.height - xAxisHeight;
            
            // Y축 범위: 현재 가격 기준 설정된 비율 (캔들의 고가/저가 고려)
            let centerPrice = currentPrice > 0 ? currentPrice : (data.length > 0 ? data[data.length - 1].close : 87000);
            if (data.length > 0) {
                // 모든 캔들의 고가와 저가를 고려하여 범위 계산
                const allHighs = data.map(c => c.high);
                const allLows = data.map(c => c.low);
                const maxHigh = Math.max(...allHighs);
                const minLow = Math.min(...allLows);
                centerPrice = (maxHigh + minLow) / 2;
            }
            const yRangePercent = window.chartConfig ? window.chartConfig.yAxisRangePercent : 0.02;
            const yRange = centerPrice * yRangePercent;
            const min = centerPrice - yRange;
            const max = centerPrice + yRange;
            const range = max - min || 0.0001;
            
            // X축: 가운데가 현재, 왼쪽은 과거, 오른쪽은 미래 (5초 단위)
            const totalPoints = maxDataPoints + futureDataPoints; // 과거 + 미래
            const currentIndex = maxDataPoints; // 가운데 지점 (현재)
            const stepX = chartWidth / (totalPoints - 1);
            
            // 그리드
            ctx.strokeStyle = window.chartConfig ? window.chartConfig.gridColor : '#333';
            ctx.lineWidth = 1;
            for (let i = 0; i <= 10; i++) {
                const y = chartY + (chartHeight / 10) * i;
                ctx.beginPath();
                ctx.moveTo(chartX, y);
                ctx.lineTo(canvas.width, y);
                ctx.stroke();
            }
            for (let i = 0; i <= 10; i++) {
                const x = chartX + (chartWidth / 10) * i;
                ctx.beginPath();
                ctx.moveTo(x, chartY);
                ctx.lineTo(x, chartY + chartHeight);
                ctx.stroke();
            }
            
            // 데이터 그리기
            if (data.length > 0) {
                
                // 캔들스틱 그리기
                if (data.length > 0) {
                    const upColor = window.chartConfig ? window.chartConfig.upColor : '#00ff00';
                    const downColor = window.chartConfig ? window.chartConfig.downColor : '#ff0000';
                    const candleWidth = (window.chartConfig ? window.chartConfig.candleWidth : 0.8) * stepX;
                    const halfWidth = candleWidth / 2;
                    
                    for (let i = 0; i < data.length; i++) {
                        const candle = data[i];
                        // 데이터가 객체인지 숫자인지 확인
                        const open = candle.open !== undefined ? candle.open : candle;
                        const high = candle.high !== undefined ? candle.high : candle;
                        const low = candle.low !== undefined ? candle.low : candle;
                        const close = candle.close !== undefined ? candle.close : candle;
                        
                        const x = chartX + (currentIndex - (data.length - 1) + i) * stepX;
                        
                        // 가격을 Y 좌표로 변환
                        const openY = chartY + chartHeight - ((open - min) / range) * chartHeight;
                        const closeY = chartY + chartHeight - ((close - min) / range) * chartHeight;
                        const highY = chartY + chartHeight - ((high - min) / range) * chartHeight;
                        const lowY = chartY + chartHeight - ((low - min) / range) * chartHeight;
                        
                        // 상승/하락 여부
                        const isUp = close >= open;
                        const color = isUp ? upColor : downColor;
                        
                        // 심지 그리기 (고가-저가)
                        ctx.strokeStyle = color;
                        ctx.lineWidth = 1;
                        ctx.beginPath();
                        ctx.moveTo(x, highY);
                        ctx.lineTo(x, lowY);
                        ctx.stroke();
                        
                        // 몸통 그리기 (시가-종가)
                        const bodyTop = Math.min(openY, closeY);
                        const bodyBottom = Math.max(openY, closeY);
                        const bodyHeight = bodyBottom - bodyTop || 1; // 최소 1px
                        
                        ctx.fillStyle = color;
                        ctx.fillRect(x - halfWidth, bodyTop, candleWidth, bodyHeight);
                        
                        // 테두리
                        ctx.strokeStyle = color;
                        ctx.strokeRect(x - halfWidth, bodyTop, candleWidth, bodyHeight);
                    }
                }
                
                // 현재 가격 점 및 수평선 표시
                if (currentPrice > 0 && data.length > 0) {
                    const currentX = chartX + currentIndex * stepX;
                    const currentY = chartY + chartHeight - ((currentPrice - min) / range) * chartHeight;
                    const lineColor = window.chartConfig ? window.chartConfig.currentPriceLineColor : '#ffffff';
                    const textColor = window.chartConfig ? window.chartConfig.currentPriceTextColor : '#ffffff';
                    
                    // 현재가 수평선 그리기
                    ctx.strokeStyle = lineColor;
                    ctx.lineWidth = 1;
                    ctx.setLineDash([5, 5]);
                    ctx.beginPath();
                    ctx.moveTo(chartX, currentY);
                    ctx.lineTo(chartWidth, currentY);
                    ctx.stroke();
                    ctx.setLineDash([]);
                    
                    // 현재 가격 점
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.currentPriceColor : '#00c853';
                    ctx.beginPath();
                    ctx.arc(currentX, currentY, 4, 0, Math.PI * 2);
                    ctx.fill();
                    
                    // 가격 텍스트 (우측에 표시, 선과 겹치지 않도록)
                    ctx.fillStyle = textColor;
                    ctx.font = '14px sans-serif';
                    ctx.textAlign = 'right';
                    // 텍스트를 선 위쪽에 표시하여 겹치지 않게
                    const textY = currentY - 15;
                    ctx.fillText(currentPrice.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2}), chartWidth - 5, textY);
                }
                
                // Y축 레이블 (금액) - 우측
                ctx.fillStyle = 'rgba(255, 255, 255, 0.6)';
                ctx.font = '12px sans-serif';
                ctx.textAlign = 'left';
                ctx.textBaseline = 'middle';
                
                for (let i = 0; i <= 10; i++) {
                    const y = chartY + (chartHeight / 10) * i;
                    const value = max - (range / 10) * i;
                    const label = value.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2});
                    ctx.fillText(label, chartWidth + 10, y);
                }
                
                // Touch/NoTouch 범위 선 그리기 (항상 표시)
                if (window.currentTradeType === 'touchnotouch' && window.touchRangeLines) {
                    // 범위 선이 없으면 초기화
                    if (window.touchRangeLines.upper === null || window.touchRangeLines.lower === null) {
                        const lastCandle = data.length > 0 ? data[data.length - 1] : null;
                        const price = currentPrice > 0 ? currentPrice : (lastCandle ? (lastCandle.close !== undefined ? lastCandle.close : (lastCandle.high !== undefined ? lastCandle.close : lastCandle)) : 87000);
                        const range = price * 0.05;
                        window.touchRangeLines.upper = price + range;
                        window.touchRangeLines.lower = price - range;
                    }
                    
                    if (window.touchRangeLines.upper !== null && window.touchRangeLines.lower !== null) {
                    const upperY = chartY + chartHeight - ((window.touchRangeLines.upper - min) / range) * chartHeight;
                    const lowerY = chartY + chartHeight - ((window.touchRangeLines.lower - min) / range) * chartHeight;
                    
                    // Y 좌표를 차트 영역 내로 제한
                    const clampedUpperY = Math.max(chartY, Math.min(chartY + chartHeight, upperY));
                    const clampedLowerY = Math.max(chartY, Math.min(chartY + chartHeight, lowerY));
                    
                    // 상단 선
                    ctx.strokeStyle = window.chartConfig ? window.chartConfig.upperRangeLineColor : '#ff1744';
                    ctx.lineWidth = 2;
                    ctx.setLineDash([5, 5]);
                    ctx.beginPath();
                    ctx.moveTo(chartX, clampedUpperY);
                    ctx.lineTo(chartWidth, clampedUpperY);
                    ctx.stroke();
                    
                    // 하단 선
                    ctx.strokeStyle = window.chartConfig ? window.chartConfig.lowerRangeLineColor : '#00e676';
                    ctx.beginPath();
                    ctx.moveTo(chartX, clampedLowerY);
                    ctx.lineTo(chartWidth, clampedLowerY);
                    ctx.stroke();
                    ctx.setLineDash([]);
                    
                    // 드래그 핸들 (원) - 우측에 표시 (더 크게 만들어서 클릭하기 쉽게)
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.upperRangeLineColor : '#ff1744';
                    ctx.beginPath();
                    ctx.arc(chartWidth - 15, clampedUpperY, 8, 0, Math.PI * 2);
                    ctx.fill();
                    
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.lowerRangeLineColor : '#00e676';
                    ctx.beginPath();
                    ctx.arc(chartWidth - 15, clampedLowerY, 8, 0, Math.PI * 2);
                    ctx.fill();
                    
                    // 가격 표시
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.upperRangeLineColor : '#ff1744';
                    ctx.font = '12px sans-serif';
                    ctx.textAlign = 'right';
                    ctx.fillText(window.touchRangeLines.upper.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2}), chartWidth - 25, clampedUpperY - 8);
                    
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.lowerRangeLineColor : '#00e676';
                    ctx.fillText(window.touchRangeLines.lower.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2}), chartWidth - 25, clampedLowerY + 18);
                    }
                }
                
                // 만기 시간 수직선 그리기 (항상 표시)
                if (expirationIndex >= 0) {
                    const expirationX = chartX + expirationIndex * stepX;
                    
                    // 수직선
                    ctx.strokeStyle = window.chartConfig ? window.chartConfig.expirationLineColor : '#ff9800';
                    ctx.lineWidth = 2;
                    ctx.setLineDash([5, 5]);
                    ctx.beginPath();
                    ctx.moveTo(expirationX, chartY);
                    ctx.lineTo(expirationX, chartY + chartHeight);
                    ctx.stroke();
                    ctx.setLineDash([]);
                    
                    // 만기 시간 표시 (미래 시간)
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.expirationLineColor : '#ff9800';
                    ctx.font = 'bold 14px sans-serif';
                    ctx.textAlign = 'center';
                    ctx.textBaseline = 'bottom';
                    
                    // 만기 시간 계산 (현재 시간 + 미래 시간)
                    const timeInterval = window.chartConfig ? window.chartConfig.timeIntervalSeconds : 5;
                    const futureSeconds = (expirationIndex - currentIndex) * timeInterval;
                    const expirationTime = new Date(Date.now() + futureSeconds * 1000);
                    const hours = expirationTime.getHours().toString().padStart(2, '0');
                    const minutes = expirationTime.getMinutes().toString().padStart(2, '0');
                    const seconds = expirationTime.getSeconds().toString().padStart(2, '0');
                    ctx.fillText(`만기: ${hours}:${minutes}:${seconds}`, expirationX, chartY - 5);
                    
                    // 만기 시간까지 남은 시간 계산
                    if (futureSeconds > 0) {
                        const min = Math.floor(futureSeconds / 60);
                        const sec = futureSeconds % 60;
                        ctx.fillText(`남은 시간: ${min}:${sec.toString().padStart(2, '0')}`, expirationX, chartY - 25);
                    }
                    
                    // 드래그 가능 표시 (손 아이콘)
                    ctx.fillStyle = window.chartConfig ? window.chartConfig.expirationLineColor : '#ff9800';
                    ctx.font = '20px sans-serif';
                    ctx.fillText('↔', expirationX, chartY + chartHeight + 20);
                }
                
                // X축 레이블 (시간) - 과거와 미래 모두 표시
                ctx.textAlign = 'center';
                ctx.textBaseline = 'top';
                ctx.fillStyle = 'rgba(255, 255, 255, 0.6)';
                ctx.font = '12px sans-serif';
                
                // 현재 시간 표시 (가운데)
                const currentX = chartX + currentIndex * stepX;
                const now = new Date();
                const nowHours = now.getHours().toString().padStart(2, '0');
                const nowMinutes = now.getMinutes().toString().padStart(2, '0');
                const nowSeconds = now.getSeconds().toString().padStart(2, '0');
                ctx.fillText(`${nowHours}:${nowMinutes}:${nowSeconds}`, currentX, chartY + chartHeight + 5);
                
                // 과거 시간 표시 (왼쪽, 설정된 시간 간격 단위)
                const timeInterval = window.chartConfig ? window.chartConfig.timeIntervalSeconds : 5;
                const pastLabels = 5;
                for (let i = 1; i <= pastLabels; i++) {
                    const pastIndex = currentIndex - i * (60 / timeInterval); // 1분 = 60초 / 시간간격
                    if (pastIndex >= 0) {
                        const x = chartX + pastIndex * stepX;
                        const pastTime = new Date(Date.now() - i * 60 * 1000);
                        const hours = pastTime.getHours().toString().padStart(2, '0');
                        const minutes = pastTime.getMinutes().toString().padStart(2, '0');
                        const seconds = pastTime.getSeconds().toString().padStart(2, '0');
                        ctx.fillText(`${hours}:${minutes}:${seconds}`, x, chartY + chartHeight + 5);
                    }
                }
                
                // 미래 시간 표시 (오른쪽, 설정된 시간 간격 단위)
                const futureLabels = 5;
                for (let i = 1; i <= futureLabels; i++) {
                    const futureIndex = currentIndex + i * (60 / timeInterval); // 1분 = 60초 / 시간간격
                    if (futureIndex < totalPoints) {
                        const x = chartX + futureIndex * stepX;
                        // 만기 시간 선과 겹치지 않도록
                        if (Math.abs(x - (chartX + expirationIndex * stepX)) > 30) {
                            const futureTime = new Date(Date.now() + i * 60 * 1000);
                            const hours = futureTime.getHours().toString().padStart(2, '0');
                            const minutes = futureTime.getMinutes().toString().padStart(2, '0');
                            const seconds = futureTime.getSeconds().toString().padStart(2, '0');
                            ctx.fillText(`${hours}:${minutes}:${seconds}`, x, chartY + chartHeight + 5);
                        }
                    }
                }
            }
        };
        
        // 마우스 이벤트 처리
        const getMousePos = (e) => {
            const rect = canvas.getBoundingClientRect();
            return {
                x: e.clientX - rect.left,
                y: e.clientY - rect.top
            };
        };
        
        const yAxisWidth = 80;
        const xAxisHeight = 30;
        
        canvas.addEventListener('mousedown', (e) => {
            const pos = getMousePos(e);
            const chartX = 0;
            const chartWidth = canvas.width - yAxisWidth;
            const chartY = 0;
            const chartHeight = canvas.height - xAxisHeight;
            
            const totalPoints = maxDataPoints + futureDataPoints;
            const currentIndex = maxDataPoints;
            const stepX = chartWidth / (totalPoints - 1);
            const expirationX = chartX + expirationIndex * stepX;
            
            // 만기 시간 선 근처 클릭 확인 (10px 범위)
            if (expirationIndex >= currentIndex && Math.abs(pos.x - expirationX) < 10 && pos.y >= chartY && pos.y <= chartY + chartHeight) {
                isDragging = true;
                canvas.style.cursor = 'grabbing';
            }
            
                // Touch/NoTouch 범위 선 드래그 확인 (차트 전체에서 드래그 가능)
                if (window.currentTradeType === 'touchnotouch' && window.touchRangeLines && window.touchRangeLines.upper !== null && window.touchRangeLines.lower !== null) {
                    const lastCandle = data.length > 0 ? data[data.length - 1] : null;
                    const centerPrice = currentPrice > 0 ? currentPrice : (lastCandle ? (lastCandle.close || lastCandle) : 87000);
                    const yRangePercent = window.chartConfig ? window.chartConfig.yAxisRangePercent : 0.02;
                    const yRange = centerPrice * yRangePercent;
                    const min = centerPrice - yRange;
                    const max = centerPrice + yRange;
                    const range = max - min || 0.0001;
                
                const upperY = chartY + chartHeight - ((window.touchRangeLines.upper - min) / range) * chartHeight;
                const lowerY = chartY + chartHeight - ((window.touchRangeLines.lower - min) / range) * chartHeight;
                
                // Y 좌표를 차트 영역 내로 제한
                const clampedUpperY = Math.max(chartY, Math.min(chartY + chartHeight, upperY));
                const clampedLowerY = Math.max(chartY, Math.min(chartY + chartHeight, lowerY));
                
                // 상단 선 또는 하단 선 근처 클릭 확인 (10px 범위, 차트 전체에서)
                if (pos.x >= chartX && pos.x <= chartWidth && pos.y >= chartY && pos.y <= chartY + chartHeight) {
                    if (Math.abs(pos.y - clampedUpperY) < 10) {
                        window.touchRangeLines.isDragging = true;
                        window.touchRangeLines.dragLine = 'upper';
                        window.touchRangeLines.dragStartY = pos.y;
                        window.touchRangeLines.dragStartUpper = window.touchRangeLines.upper;
                        window.touchRangeLines.dragStartLower = window.touchRangeLines.lower;
                        window.touchRangeLines.dragStartCenter = (window.touchRangeLines.upper + window.touchRangeLines.lower) / 2;
                        window.touchRangeLines.dragStartRange = window.touchRangeLines.upper - window.touchRangeLines.lower;
                        canvas.style.cursor = 'ns-resize';
                        e.preventDefault();
                    } else if (Math.abs(pos.y - clampedLowerY) < 10) {
                        window.touchRangeLines.isDragging = true;
                        window.touchRangeLines.dragLine = 'lower';
                        window.touchRangeLines.dragStartY = pos.y;
                        window.touchRangeLines.dragStartUpper = window.touchRangeLines.upper;
                        window.touchRangeLines.dragStartLower = window.touchRangeLines.lower;
                        window.touchRangeLines.dragStartCenter = (window.touchRangeLines.upper + window.touchRangeLines.lower) / 2;
                        window.touchRangeLines.dragStartRange = window.touchRangeLines.upper - window.touchRangeLines.lower;
                        canvas.style.cursor = 'ns-resize';
                        e.preventDefault();
                    }
                }
            }
        });
        
        canvas.addEventListener('mousemove', (e) => {
            const pos = getMousePos(e);
            const chartX = 0;
            const chartWidth = canvas.width - yAxisWidth;
            const chartY = 0;
            const chartHeight = canvas.height - xAxisHeight;
            
            if (isDragging) {
                const totalPoints = maxDataPoints + futureDataPoints;
                const currentIndex = maxDataPoints;
                const stepX = chartWidth / (totalPoints - 1);
                // 마우스 위치에 가장 가까운 인덱스 계산 (5초 단위)
                const newIndex = Math.round((pos.x - chartX) / stepX);
                // 현재 인덱스 이후(미래)로만 이동 가능
                expirationIndex = Math.max(currentIndex, Math.min(newIndex, totalPoints - 1));
                drawChart();
            } else if (window.touchRangeLines && window.touchRangeLines.isDragging && window.currentTradeType === 'touchnotouch') {
                const centerPrice = currentPrice > 0 ? currentPrice : (data.length > 0 ? data[data.length - 1] : 87000);
                const yRangePercent = window.chartConfig ? window.chartConfig.yAxisRangePercent : 0.02;
                const yRange = centerPrice * yRangePercent;
                const min = centerPrice - yRange;
                const max = centerPrice + yRange;
                const range = max - min || 0.0001;
                
                // Y 좌표를 가격으로 변환
                const newPrice = max - ((pos.y - chartY) / chartHeight) * range;
                
                // 시작 중심점과 범위
                const startCenter = window.touchRangeLines.dragStartCenter;
                const startRange = window.touchRangeLines.dragStartRange;
                
                if (window.touchRangeLines.dragLine === 'upper') {
                    // 상단 선 드래그: 중심점 기준으로 같은 비율로 위아래 범위 조정
                    const deltaY = pos.y - window.touchRangeLines.dragStartY;
                    const deltaPrice = -((deltaY / chartHeight) * range); // Y는 위에서 아래로 증가하므로 반대
                    
                    // 새로운 상단 가격
                    const newUpper = window.touchRangeLines.dragStartUpper + deltaPrice;
                    
                    // 중심점은 유지하고, 범위를 조정 (위아래 같은 비율로)
                    const center = startCenter;
                    const newRange = (newUpper - center) * 2; // 중심점에서 위로의 거리 * 2
                    
                    window.touchRangeLines.upper = center + newRange / 2;
                    window.touchRangeLines.lower = center - newRange / 2;
                    
                    // 범위 체크
                    if (window.touchRangeLines.lower < min) {
                        window.touchRangeLines.lower = min;
                        window.touchRangeLines.upper = min + newRange;
                    }
                    if (window.touchRangeLines.upper > max) {
                        window.touchRangeLines.upper = max;
                        window.touchRangeLines.lower = max - newRange;
                    }
                } else if (window.touchRangeLines.dragLine === 'lower') {
                    // 하단 선 드래그: 중심점 기준으로 같은 비율로 위아래 범위 조정
                    const deltaY = pos.y - window.touchRangeLines.dragStartY;
                    const deltaPrice = -((deltaY / chartHeight) * range); // Y는 위에서 아래로 증가하므로 반대
                    
                    // 새로운 하단 가격
                    const newLower = window.touchRangeLines.dragStartLower + deltaPrice;
                    
                    // 중심점은 유지하고, 범위를 조정 (위아래 같은 비율로)
                    const center = startCenter;
                    const newRange = (center - newLower) * 2; // 중심점에서 아래로의 거리 * 2
                    
                    window.touchRangeLines.upper = center + newRange / 2;
                    window.touchRangeLines.lower = center - newRange / 2;
                    
                    // 범위 체크
                    if (window.touchRangeLines.upper > max) {
                        window.touchRangeLines.upper = max;
                        window.touchRangeLines.lower = max - newRange;
                    }
                    if (window.touchRangeLines.lower < min) {
                        window.touchRangeLines.lower = min;
                        window.touchRangeLines.upper = min + newRange;
                    }
                }
                
                drawChart();
            } else if (data.length > 1) {
                const stepX = chartWidth / (data.length - 1);
                const expirationX = chartX + expirationIndex * stepX;
                
                let cursor = 'default';
                
                // 만기 시간 선 근처에 마우스가 있으면 커서 변경
                if (Math.abs(pos.x - expirationX) < 10 && pos.y > 0 && pos.y < canvas.height - xAxisHeight) {
                    cursor = 'grab';
                }
                
                // Touch/NoTouch 범위 선 근처 확인 (차트 전체에서)
                if (window.currentTradeType === 'touchnotouch' && window.touchRangeLines && window.touchRangeLines.upper !== null && window.touchRangeLines.lower !== null) {
                    const min = Math.min(...data);
                    const max = Math.max(...data);
                    const range = max - min || 0.0001;
                    
                    const upperY = chartY + chartHeight - ((window.touchRangeLines.upper - min) / range) * chartHeight;
                    const lowerY = chartY + chartHeight - ((window.touchRangeLines.lower - min) / range) * chartHeight;
                    
                    // Y 좌표를 차트 영역 내로 제한
                    const clampedUpperY = Math.max(chartY, Math.min(chartY + chartHeight, upperY));
                    const clampedLowerY = Math.max(chartY, Math.min(chartY + chartHeight, lowerY));
                    
                    // 차트 전체에서 선 근처 확인
                    if (pos.x >= chartX && pos.x <= chartWidth && pos.y >= chartY && pos.y <= chartY + chartHeight) {
                        if (Math.abs(pos.y - clampedUpperY) < 10 || Math.abs(pos.y - clampedLowerY) < 10) {
                            cursor = 'ns-resize';
                        }
                    }
                }
                
                canvas.style.cursor = cursor;
            }
        });
        
        canvas.addEventListener('mouseup', () => {
            isDragging = false;
            if (window.touchRangeLines) {
                window.touchRangeLines.isDragging = false;
                window.touchRangeLines.dragLine = null;
            }
            canvas.style.cursor = 'default';
        });
        
        canvas.addEventListener('mouseleave', () => {
            isDragging = false;
            if (window.touchRangeLines) {
                window.touchRangeLines.isDragging = false;
                window.touchRangeLines.dragLine = null;
            }
            canvas.style.cursor = 'default';
        });
        
        // drawChart를 전역으로 노출 (선택기에서 호출하기 위해)
        window.drawChart = drawChart;
        
        // 초기 과거 데이터 로드
        fetchHistoricalData();
        
        // 1초마다 가격 업데이트
        setInterval(fetchPrice, 1000);
        
        // 리사이즈
        window.addEventListener('resize', () => {
            canvas.width = chartContainer.clientWidth;
            canvas.height = chartContainer.clientHeight;
            drawChart();
        });
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const controller = new MainPageController();
    controller.pageInit();
});
