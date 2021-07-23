package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("html").Parse(tplString)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(os.Stdout, string(data))
	if err != nil {
		log.Fatal(err)
	}
}

var tplString = `
<!DOCTYPE html>

<head>
    <meta charset="utf-8">
    <style>
        body {
            font: 10px sans-serif;
        }

        text {
            fill: #000;
        }

        button {
            position: absolute;
            right: 20px;
            top: 440px;
            display: none;
        }

        path.candle {
            stroke: #000000;
        }

        path.candle.body {
            stroke-width: 0;
        }

        path.candle.up {
            fill: #00AA00;
            stroke: #00AA00;
        }

        path.candle.down {
            fill: #FF0000;
            stroke: #FF0000;
        }

        path.tradearrow {
            stroke: none;
        }

        path.tradearrow.buy {
            fill: #009900;
        }

        path.tradearrow.buy-pending {
            fill-opacity: 0.2;
            stroke: #009900;
            stroke-width: 1.5;
        }

        path.tradearrow.sell {
            fill: #990000;
        }

        .tradearrow path.highlight {
            fill: none;
            stroke-width: 2;
        }

        .tradearrow path.highlight.buy,
        .tradearrow path.highlight.buy-pending {
            stroke: #009900;
        }

        .tradearrow path.highlight.buy-pending {
            fill: #009900;
            fill-opacity: 0.3;
        }

        .tradearrow path.highlight.sell {
            stroke: #9900FF;
        }
    </style>
    <style>
        .header_section {
            font-size: 16px;
        }

        .label {
            font-size: 20px;
            margin-top: 10px;
            margin-bottom: 5px;
            font-weight: bolder;
        }
    </style>
</head>

<body>
    <button>Update</button>
    <script src="https://d3js.org/d3.v4.min.js"></script>
    <script src="http://techanjs.org/techan.min.js"></script>
    <script>
        var output = {{.}}
	
	</script>

	<section class="header_section">
		<div class="label">Exchange</div>
		<div id="exchange"></div>
		<div class="label">Symbol</div>
		<div id="symbol"></div>
		<div class="label">Take Profit Ratio</div>
		<div id="takeProfitRatio"></div>
		<div id=""></div>
		<div id=""></div>
		<div id=""></div>
		<div id=""></div>
		<div id=""></div>
		<div id="chart"></div>
	</section>

	<script>
		document.querySelector('#exchange').innerHTML = output.input.exchange
		document.querySelector('#symbol').innerHTML = output.input.baseAsset+'/'+output.input.quoteAsset
		document.querySelector('#takeProfitRatio').innerHTML = (output.profitRatio * 100.0).toFixed(2)+'%'

		var margin = { top: 20, right: 20, bottom: 30, left: 50 },
			width = 960 - margin.left - margin.right,
			height = 500 - margin.top - margin.bottom;

		var dateFormat = d3.timeFormat('` + "%s" + `'),
			parseDate = d3.timeParse('` + "%s" + `'),
			parseISODate = d3.timeParse('` + "%Y-%m-%dT%H:%M:%SZ" + `'),
			valueFormat = d3.format(',.2f');

		var x = techan.scale.financetime()
			.range([0, width]);

		var y = d3.scaleLinear()
			.range([height, 0]);

		var candlestick = techan.plot.candlestick()
			.xScale(x)
			.yScale(y);

		var tradearrow = techan.plot.tradearrow()
			.xScale(x)
			.yScale(y)
			.orient(function (d) { return d.type.startsWith("buy") ? "up" : "down"; })
			.on("mouseenter", enter)
			.on("mouseout", out);

		var xAxis = d3.axisBottom(x);

		var yAxis = d3.axisLeft(y);

		var svg = d3.select("#chart").append("svg")
			.attr("width", width + margin.left + margin.right)
			.attr("height", height + margin.top + margin.bottom)
			.append("g")
			.attr("transform", "translate(" + margin.left + "," + margin.top + ")");

		var valueText = svg.append('text')
			.style("text-anchor", "end")
			.attr("class", "coords")
			.attr("x", width - 5)
			.attr("y", 15);

		var accessor = candlestick.accessor();

		function draw(data, trades) {
			x.domain(data.map(candlestick.accessor().d));
			y.domain(techan.scale.plot.ohlc(data, candlestick.accessor()).domain());

			svg.selectAll("g.candlestick").datum(data).call(candlestick);
			svg.selectAll("g.tradearrow").datum(trades).call(tradearrow);

			svg.selectAll("g.x.axis").call(xAxis);
			svg.selectAll("g.y.axis").call(yAxis);
		}

		function enter(d) {
			valueText.style("display", "inline");
			refreshText(d);
		}

		function out() {
			valueText.style("display", "none");
		}

		function refreshText(d) {
			valueText.text("Trade: " + dateFormat(d.date) + ", " + d.type + ", " + valueFormat(d.price));
		}

		data = output.candlesticks.map(function (d) {
			return {
				date: parseDate(d.t),
				open: +d.o,
				high: +d.h,
				low: +d.l,
				close: +d.c,
				volume: +d.v
			};
		}).sort(function (a, b) { return d3.ascending(accessor.d(a), accessor.d(b)); });

		var trades = output.events.map((e) => ({
			date: parseISODate(e.at),
			type: e.eventType == 'entered' ? 'buy' : 'sell',
			price: e.price,
			quantity: 1
		}))
		// var trades = [
		//     { date: data[67].date, type: "buy", price: data[67].low, quantity: 1000 },
		//     { date: data[100].date, type: "sell", price: data[100].high, quantity: 200 },
		//     { date: data[156].date, type: "buy", price: data[156].open, quantity: 500 },
		//     { date: data[167].date, type: "sell", price: data[167].close, quantity: 300 },
		//     { date: data[187].date, type: "buy-pending", price: data[187].low, quantity: 300 }
		// ];

		svg.append("g")
			.attr("class", "candlestick");

		svg.append("g")
			.attr("class", "tradearrow");

		svg.append("g")
			.attr("class", "x axis")
			.attr("transform", "translate(0," + height + ")");

		svg.append("g")
			.attr("class", "y axis")
			.append("text")
			.attr("transform", "rotate(-90)")
			.attr("y", 6)
			.attr("dy", ".71em")
			.style("text-anchor", "end")
			.text("Price ($)");

		// Data to display initially
		draw(data, trades);
		// Only want this button to be active if the data has loaded
		// d3.select("button").on("click", function () { draw(data, trades); }).style("display", "inline");

	</script>

</body>

</html>
`
