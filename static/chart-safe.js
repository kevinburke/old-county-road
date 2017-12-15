$(function() {
  const curbHeight = 5.7;
  const laneWidth = 18.74;
  const turnoutWidth = 23.78;
  const lineWidth = 0.5;
  const curbPosition = -2;
  const bounds = 220;
  var d2 = [
    // south edge of curb
    [-bounds, curbPosition], [-39, curbPosition], null, [11, curbPosition], [bounds, curbPosition], null,
    // north edge of curb
    [-bounds, curbPosition + curbHeight], [bounds, curbPosition + curbHeight], null,
    // far side of the road
    [-bounds, curbPosition + curbHeight + laneWidth*2], [bounds, curbPosition + curbHeight + laneWidth*2], null,

    // Vertical turnout lines
    [-39, curbPosition], [-39, -100], null,
    [-39 + turnoutWidth, curbPosition], [-39 + turnoutWidth, -100], null,
    [11 - turnoutWidth, curbPosition], [11 - turnoutWidth, -100], null,
    [11, curbPosition], [11, -100], null,

    // Curb section where two turnouts meet
    [-39 + turnoutWidth, curbPosition], [11 - turnoutWidth, curbPosition], null,

  ];

  var dashedLane = [
    [-bounds, curbPosition + curbHeight + laneWidth], [bounds, curbPosition + curbHeight + laneWidth],
  ];
  var redCurb = [
    [-39, curbPosition + curbHeight + 0.3], [-39 -96.62, curbPosition + curbHeight + 0.3], null,
    [11, curbPosition + curbHeight + 0.3], [11 + 54.5, curbPosition + curbHeight + 0.3], null,
    [-200, curbPosition - 10], [-200, curbPosition + curbHeight + laneWidth * 2 + 10], null,
    [200, curbPosition - 10], [200, curbPosition + curbHeight + laneWidth * 2 + 10], null,
  ];

  const leftCurbCarStart = -39 - 96.62 - 0.5;
  const leftCurbCarEnd = leftCurbCarStart - 14.288;
  const rightCurbCarStart = 11 + 54.5 + 0.5;
  const rightCurbCarEnd = rightCurbCarStart + 14.288;
  var curbCar = [
    [leftCurbCarStart, curbPosition + curbHeight + 0.5], [leftCurbCarStart, curbPosition+curbHeight + 0.5 + 5.75], [leftCurbCarEnd, curbPosition+curbHeight + 0.5 + 5.75], [leftCurbCarEnd, curbPosition + curbHeight + 0.5], [leftCurbCarStart, curbPosition + curbHeight + 0.5], null,
    [rightCurbCarStart, curbPosition + curbHeight + 0.5], [rightCurbCarStart, curbPosition+curbHeight + 0.5 + 5.75], [rightCurbCarEnd, curbPosition+curbHeight + 0.5 + 5.75], [rightCurbCarEnd, curbPosition + curbHeight + 0.5], [rightCurbCarStart, curbPosition + curbHeight + 0.5], null,
  ];

  var sightLines = [
    //[0, 0], [leftCurbCarStart, curbPosition + curbHeight + 5.75], null,
    [0, 0], [-bounds, 14.86], null,

    // RHS
    //[0, 0], [rightCurbCarStart, curbPosition + curbHeight + 5.75], null,
    [0, 0], [bounds, 31.5], null,

  ];

  $.plot("#placeholder-2", [ {
    data: d2,
    lines: {show: true, lineWidth: lineWidth},
    color: "#333",
    shadowSize: 1,
  }, {
    data: curbCar,
    lines: {show: true, lineWidth: lineWidth},
    color: "#333",
    shadowSize: 1,
  }, {
    data: dashedLane,
    lines: {show: true, lineWidth: lineWidth},
    dashes: {show: true},
    shadowSize: 1,
  }, {
    data: redCurb,
    lines: {show: true, lineWidth: 2},
    color: "#d14836",
    shadowSize: 1,
  }, {
    data: sightLines,
    lines: {show: true, lineWidth: lineWidth},
    color: "#1069da",
    shadowSize: 1,
  },
  ], {grid: {show: false}});
});
