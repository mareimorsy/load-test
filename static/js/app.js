// $.get('/load', function(data){





  // set the dimensions and margins of the graph
  var margin = {top: 10, right: 30, bottom: 30, left: 60},
  width = 460 - margin.left - margin.right,
  height = 400 - margin.top - margin.bottom;

  // append the svg object to the body of the page
  var svg = d3.select("#my_dataviz")
  .append("svg")
  .attr("width", width + margin.left + margin.right)
  .attr("height", height + margin.top + margin.bottom)
  .append("g")
  .attr("transform",
        "translate(" + margin.left + "," + margin.top + ")");

  //Read the data
  d3.json("/load", function(data) {
    // console.log(data[0].start)

    var reqArr = data

    // reqArr = JSON.parse(data)
    var timeStart = 0
    var timeEnd = reqArr[0].start || 0
  
    var minLatency = reqArr[0].latency
    var maxLatency = reqArr[0].latency
  
    for (i = 0; i< reqArr.length; i++){
      if (reqArr[i].latency < minLatency)
        minLatency = reqArr[i].latency
      if (reqArr[i].latency > maxLatency)
        maxLatency = reqArr[i].latency
      if (reqArr[i].start > timeEnd)
        timeEnd = reqArr[i].start
    }

  // Add X axis
  var x = d3.scaleLinear()
  .domain([timeStart, timeEnd])
  .range([ 0, width ]);
  svg.append("g")
  .attr("transform", "translate(0," + height + ")")
  .call(d3.axisBottom(x));

  // Add Y axis
  var y = d3.scaleLinear()
  .domain([minLatency, maxLatency])
  .range([ height, 0]);
  svg.append("g")
  .call(d3.axisLeft(y));

  // Add dots
  svg.append('g')
  .selectAll("dot")
  .data(data)
  .enter()
  .append("circle")
    .attr("cx", function (d) { return x(d.start); } )
    .attr("cy", function (d) { return y(d.latency); } )
    .attr("r", 1.5)
    .style("fill", "#69b3a2")

  })
// })



