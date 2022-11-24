## APIs
#### Methods
```
[
        '휴식', // 0
        'Kizen Powerlifting Peaking Program', // 1
        'nSuns Programs',
        'Jim Wendler 5/3/1 Programs',
        'Calgary Barbell Programs',
        'Sheiko Programs',
        'Candito Program',
        'Juggernaut Method Base Template',
        'Greg Nuckols 28 Programs',
        'Beginner Powerlifting Programs',
        'Intermediate Powerlifting Programs',
        'Madcow 5x5 Program',
        'General 5x5 Program',
        '기타(etc)' // 기타는 index 242로 설정 [ 주의 ]
    ]
```
#### POST api/pred
```
// Request body
{
  recordFixed: {
    sex: string; // male|female
    age: Number;
    height: Number;
    arm_length: Number; // long | medium | short
    leg_length: Number; // long | medium | short
  },
  recordWeekly: [
    {
      weight: Number;
      fat_rate: Number;
      program: Number; // 그 주에 한 프로그램 인덱스
      squat: Number; // 그 주의 프로그램을 끝내고 스쿼트, 벤치, 데드 무게
      benchpress: Number;
      deadlift: Number;
    }, ... // 24 weeks
  ]
}

// Response
[
  { method: 'Wendler 531', sqaut: 162.3, benchpress: 120.1, deadlift: 181.5 },
  ...
  // best 5 methods and results  
```
response에서 최고의 5개 메소드를 받아서 이 훈련 프로그램을 사용하면 이 정도 증량을 얻을거라 보여주면됨.
