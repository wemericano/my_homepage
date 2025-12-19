// 큰제목별 작은제목 목록
const subCategories = {
    'IT': [
        '404 에러', '서버 다운', '네트워크 장애', '보안 취약점', '데이터 백업',
        '시스템 업그레이드', '클라우드 마이그레이션', 'IT 인프라', '디지털 트랜스포메이션'
    ],
    '개발': [
        '코드 리뷰', '버전 관리', '테스트 자동화', 'CI/CD', '코드 품질',
        '리팩토링', '디버깅', '성능 최적화', '아키텍처 설계', '개발 프로세스'
    ],
    '프로그래밍': [
        '변수와 함수', '객체지향 프로그래밍', '알고리즘', '자료구조', '디자인 패턴',
        '예외 처리', '메모리 관리', '동시성 프로그래밍', '함수형 프로그래밍'
    ],
    '웹개발': [
        'HTML/CSS', 'JavaScript', 'React', 'Vue.js', 'Angular', 'Node.js',
        'RESTful API', '웹 성능 최적화', '반응형 디자인', 'PWA', '웹 보안'
    ],
    '모바일': [
        'iOS 개발', 'Android 개발', 'React Native', 'Flutter', '모바일 UI/UX',
        '앱 성능 최적화', '푸시 알림', '인앱 결제', '모바일 보안'
    ],
    '인공지능': [
        '머신러닝', '딥러닝', '자연어 처리', '컴퓨터 비전', '강화학습',
        '신경망', '데이터 전처리', '모델 학습', 'AI 윤리'
    ],
    '보안': [
        'SQL Injection', 'XSS 공격', 'CSRF', '암호화', '인증/인가',
        '보안 취약점 분석', '침입 탐지', '화이트햇 해킹', '보안 정책'
    ],
    '데이터베이스': [
        'SQL 쿼리', '인덱싱', '트랜잭션', '정규화', 'NoSQL',
        '데이터베이스 설계', '성능 튜닝', '백업/복구', '데이터 모델링'
    ],
    '클라우드': [
        'AWS', 'Azure', 'GCP', '컨테이너', '쿠버네티스',
        '서버리스', '클라우드 마이그레이션', '클라우드 보안', '비용 최적화'
    ],
    'DevOps': [
        'Docker', 'Kubernetes', 'Jenkins', 'GitLab CI', '모니터링',
        '로깅', '인프라 자동화', '마이크로서비스', '배포 전략'
    ],
    '프론트엔드': [
        'React Hooks', '상태 관리', '컴포넌트 설계', '웹 접근성', '브라우저 호환성',
        '번들러', '타입스크립트', 'CSS 프레임워크', '프론트엔드 테스트'
    ],
    '백엔드': [
        'API 설계', '서버 아키텍처', '인증 시스템', '캐싱', '메시지 큐',
        '마이크로서비스', '서버 성능', '데이터베이스 연결', '로드 밸런싱'
    ],
    '풀스택': [
        '전체 시스템 설계', '프론트엔드-백엔드 연동', '풀스택 프레임워크',
        '데이터 흐름', '전체 개발 프로세스', '풀스택 개발자 역량'
    ],
    '알고리즘': [
        '정렬 알고리즘', '탐색 알고리즘', '동적 프로그래밍', '그래프 알고리즘',
        '시간 복잡도', '공간 복잡도', '알고리즘 최적화', '코딩 테스트'
    ],
    '네트워크': [
        'HTTP/HTTPS', 'TCP/IP', 'DNS', '로드 밸런싱', 'CDN',
        '네트워크 보안', '프로토콜', '네트워크 성능', 'API 게이트웨이'
    ]
};

document.addEventListener('DOMContentLoaded', function() {
    const mainCategory = document.getElementById('mainCategory');
    const subCategory = document.getElementById('subCategory');
    const generateBtn = document.getElementById('generateBtn');

    // 큰제목 선택 시 작은제목 목록 업데이트
    mainCategory.addEventListener('change', function() {
        const selectedMain = this.value;
        subCategory.innerHTML = '<option value="">선택하세요</option>';
        
        if (selectedMain && subCategories[selectedMain]) {
            subCategories[selectedMain].forEach(sub => {
                const option = document.createElement('option');
                option.value = sub;
                option.textContent = sub;
                subCategory.appendChild(option);
            });
        }
    });

    // 생성 버튼 클릭 이벤트
    generateBtn.addEventListener('click', function() {
        const main = mainCategory.value;
        const sub = subCategory.value;

        if (!main || !sub) {
            alert('큰제목과 작은제목을 모두 선택해주세요.');
            return;
        }

        // 로딩 상태 시작
        generateBtn.disabled = true;
        generateBtn.innerHTML = '<span class="loading-spinner"></span>생성 중...';
        const loadingText = document.getElementById('loadingText');
        loadingText.classList.add('show');
        
        // 결과 영역 숨기기
        const resultArea = document.getElementById('resultArea');
        resultArea.classList.remove('show');

        // 서버로 요청 보내기
        fetch('/api/blog/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                mainCategory: main,
                subCategory: sub
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('서버 요청 실패');
            }
            return response.json();
        })
        .then(data => {
            console.log('서버 응답:', data);
            
            // 로딩 상태 종료
            generateBtn.disabled = false;
            generateBtn.innerHTML = '생성';
            loadingText.classList.remove('show');
            
            // 결과 영역 표시
            const resultArea = document.getElementById('resultArea');
            const resultMsg = document.getElementById('resultMsg');
            const resultTitle = document.getElementById('resultTitle');
            const resultContent = document.getElementById('resultContent');
            
            resultMsg.textContent = data.message || '-';
            resultTitle.textContent = data.data?.title || '-';
            resultContent.textContent = data.data?.content || '-';
            
            resultArea.classList.add('show');
            
            // 결과 영역으로 스크롤
            resultArea.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        })
        .catch(error => {
            console.error('에러:', error);
            
            // 로딩 상태 종료
            generateBtn.disabled = false;
            generateBtn.innerHTML = '생성';
            loadingText.classList.remove('show');
            
            alert('서버 요청 중 오류가 발생했습니다.');
        });
    });
});

