$(function() {
  const curbHeight = 5.7;
  const laneWidth = 18.74;
  const turnoutWidth = 23.78;
  const lineWidth = 0.5;
  const curbPosition = -2;
  var d2 = [
    // south edge of curb
    [-220, curbPosition], [-45, curbPosition], null, [5, curbPosition], [220, curbPosition], null,
    // north edge of curb
    [-220, curbPosition + curbHeight], [220, curbPosition + curbHeight], null,
    // far side of the road
    [-220, curbPosition + curbHeight + laneWidth*2], [220, curbPosition + curbHeight + laneWidth*2], null,

    // Vertical turnout lines
    [-45, curbPosition], [-45, -100], null,
    [-45 + turnoutWidth, curbPosition], [-45 + turnoutWidth, -100], null,
    [5 - turnoutWidth, curbPosition], [5 - turnoutWidth, -100], null,
    [5, curbPosition], [5, -100], null,

    // Curb section where two turnouts meet
    [-45 + turnoutWidth, curbPosition], [5 - turnoutWidth, curbPosition], null,

  ];

  var dashedLane = [
    [-220, curbPosition + curbHeight + laneWidth], [220, curbPosition + curbHeight + laneWidth],
  ];
  var redCurb = [
    [-45, curbPosition + curbHeight + 0.3], [-68.05, curbPosition + curbHeight + 0.3], null,
    [-200, curbPosition - 10], [-200, curbPosition + curbHeight + laneWidth * 2 + 10], null,
    [200, curbPosition - 10], [200, curbPosition + curbHeight + laneWidth * 2 + 10], null,
  ];

  var curbCar = [
    [-68.05 - 0.5, curbPosition + curbHeight + 0.5], [-68.05 - 0.5, curbPosition+curbHeight + 0.5 + 5.75], [-68.05 - 0.5 - 14.288, curbPosition+curbHeight + 0.5 + 5.75], [-68.05 - 0.5 - 14.288, curbPosition + curbHeight + 0.5], [-68.05 - 0.5, curbPosition + curbHeight + 0.5], null,

    [5 + 0.5, curbPosition + curbHeight + 0.5], [5 + 0.5, curbPosition+curbHeight + 0.5 + 5.75], [5 + 0.5 + 14.288, curbPosition+curbHeight + 0.5 + 5.75], [5 + 0.5 + 14.288, curbPosition + curbHeight + 0.5], [5 + 0.5, curbPosition + curbHeight + 0.5], null,
  ];

  var sightLines = [
    // [0, 0], [-52*3, curbPosition + curbHeight + laneWidth - 2], null,
    [0, 0], [-171.1878, curbPosition + curbHeight + laneWidth], null,

    // RHS
    //[0, 0], [23, curbPosition + curbHeight + laneWidth + 12], null,
    [0, 0], [27.50116, curbPosition + curbHeight + laneWidth*2], null,

    [0, 0], [229.69, curbPosition + curbHeight + laneWidth + 12], null,
  ];

  $.plot("#placeholder", [ {
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
