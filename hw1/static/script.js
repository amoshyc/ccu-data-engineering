var Signal = function () {
    var obj = {
        slots: [],
        register(slot) {
            obj.slots.push(slot);
        },
        trigger(data) {
            for (var idx in obj.slots) {
                obj.slots[idx](data);
            }
        }
    };
    return obj;
};

var signals = {
    fab_clk: Signal(), // fab clicked by user
    res_rec: Signal(), // search result received
    stat_clk: Signal(), // stat button clicked by user
    mark_clk: Signal(), // mark button clicked by user
    file_clk: Signal(), // info button clicked by user
    upd_stat: Signal(), // update output to mode stat
    upd_mark: Signal(), // update output to mode mark
    upd_file: Signal(), // update output to mode file
};

var output = {
    elem: $('#output'),
    stat_elem: $('#stat'),
    mark_elem: $('#mark'),
    file_elem: $('#file'),
    chart: null,
    init: function () {
        signals['upd_stat'].register(output.update_stat);
        signals['upd_mark'].register(output.update_mark);
        signals['upd_file'].register(output.update_file);
    },
    init_chart: function () {
        Chart.defaults.global.defaultFontFamily = 'Bitter, sans';
        output.chart = new Chart(output.stat_elem, {
            type: 'horizontalBar',
            data: {
                labels: [],
                datasets: [{
                    label: '# of Occurrence',
                    data: [],
                    borderWidth: 1
                }]
            },
            options: {
                title: {
                    display: true,
                    text: 'Top 10 terms'
                },
                legend: {
                    display: false
                },
                tooltips: {
                    displayColors: false
                },
                scales: {
                    xAxes: [{
                        position: 'top',
                        ticks: {
                            suggestedMin: 0,
                            suggestedMax: 10,
                            beginAtZero: true
                        }
                    }]
                }
            }
        });
    },
    update_stat: function (args) {
        if (output.chart == null) {
            output.init_chart();
        }

        var [terms, counts] = args;
        if (terms.length > 10) {
            terms = terms.slice(0, 10);
            counts = counts.slice(0, 10);
        }

        output.chart.data.labels = terms;
        output.chart.data.datasets[0].data = counts;
        output.chart.update();

        output.file_elem.hide();
        output.mark_elem.hide();
        output.stat_elem.show();
    },
    update_mark: function (data) {
        var lines = $('#input2').val().split('\n');
        var all_pos = {};
        for (var i = 0; i < lines.length; i++) {
            all_pos[i] = [];
        }
        for (var i in data) {
            var term = data[i][0];
            var term_pos = data[i][1];
            for (var j in term_pos) {
                var row = term_pos[j][0];
                var col = term_pos[j][1];
                all_pos[row].push([col, col + term.length]);
            }
        }
        for (var i = 0; i < lines.length; i++) {
            if (all_pos[i].length > 0) {
                all_pos[i].sort(function (a, b) {
                    if (a[0] == b[0]) {
                        return a[1] - b[1];
                    }
                    return a[0] - b[0];
                });
            }
        }

        var html = '';
        for (var i in lines) {
            var line = lines[i];
            if (all_pos[i].length == 0) {
                html += line;
            }
            else {
                var idx = 0;
                for (var j in all_pos[i]) {
                    var [s, t] = all_pos[i][j];
                    if (s > idx) {
                        html += line.slice(idx, s);
                    }
                    if (t > s) {
                        html += '<mark>' + line.slice(s, t) + '</mark>';
                    }
                    idx = t;
                }
                if (idx < line.length) {
                    html += line.slice(idx, line.length);
                }
            }
            html += '<br>';
        }

        output.mark_elem.html(html);
        output.stat_elem.hide();
        output.file_elem.hide();
        output.mark_elem.show();
    },
    update_file: function(file) {
        output.file_elem.html(file);

        output.stat_elem.hide();
        output.mark_elem.hide();
        output.file_elem.show();
    }
};

var view = {
    elem_stat: $('#btn-stat'),
    elem_mark: $('#btn-mark'),
    elem_file: $('#btn-file'),
    default_color: '#777',
    active_color: '#2196f3',
    init: function () {
        view.elem_stat.click(function () {
            signals['stat_clk'].trigger();
        });
        view.elem_mark.click(function () {
            signals['mark_clk'].trigger();
        });
        view.elem_file.click(function() {
            signals['file_clk'].trigger();
        })
        signals['upd_stat'].register(view.update_button);
        signals['upd_mark'].register(view.update_button);
        signals['upd_file'].register(view.update_button);
    },
    update_button: function () {
        if (model.output_mode == 'stat') {
            view.elem_stat.css('color', view.active_color);
            view.elem_mark.css('color', view.default_color);
            view.elem_file.css('color', view.default_color);            
        }
        else if (model.output_mode == 'mark') {
            view.elem_stat.css('color', view.default_color);
            view.elem_mark.css('color', view.active_color);
            view.elem_file.css('color', view.default_color);
            
        }
        else {
            view.elem_stat.css('color', view.default_color);
            view.elem_mark.css('color', view.default_color);
            view.elem_file.css('color', view.active_color);
        }
    }
};

var time = {
    elem: $('#time'),
    init: function () {
        signals['fab_clk'].register(function () {
            time.elem.html('Running');
        });
        signals['res_rec'].register(function (res) {
            time.elem.hide();
            time.elem.html(res['time']);
            time.elem.fadeIn('fast');
        });
    }
}

var fab = {
    elem: $('#fab'),
    init: function () {
        signals['fab_clk'].register(fab.onclick);
        fab.elem.click(function () {
            signals['fab_clk'].trigger();
        });
    },
    onclick() {
        data = {
            'term': $('#input1').val(),
            'text': $('#input2').val(),
            'acc': $('#switch-1').prop('checked'),
            'use': $('#switch-2').prop('checked'),
            'mp': $('#switch-3').prop('checked')
        };
        $.post('./data', data, function (ret) {
            signals['res_rec'].trigger(ret);
        });
    }
};


var model = {
    data: null,
    file: null,
    output_mode: 'stat', // 'stat', 'mark', or 'file
    init: function () {
        signals['res_rec'].register(model.process_result);
        signals['stat_clk'].register(function () {
            model.change_mode('stat');
        });
        signals['mark_clk'].register(function () {
            model.change_mode('mark');
        });
        signals['file_clk'].register(function () {
            model.change_mode('file');
        });
    },
    get_terms_counts: function () {
        var terms = [];
        var counts = [];
        for (var i in model.data) {
            terms.push(model.data[i][0]);
            counts.push(model.data[i][1].length);
        }
        return [terms, counts];
    },
    process_result: function (ret) {
        console.log(ret);
        model.file = ret['file'];
        model.data = ret['res'];
        if (ret['file_only']) {
            model.change_mode('file');
        }
        else {
            model.change_mode(model.output_mode);
        }
    },
    change_mode: function (new_mode) {
        model.output_mode = new_mode;
        if (new_mode == 'mark') {
            signals['upd_mark'].trigger(model.data);
        }
        else if (new_mode == 'stat') {
            signals['upd_stat'].trigger(model.get_terms_counts());
        }
        else {
            signals['upd_file'].trigger(model.file);
        }
    }
};

$(function () {
    // view
    fab.init();
    view.init();
    output.init();
    time.init();
    // model
    model.init();
});