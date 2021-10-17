$(() => {
    const model = {};
    const update = () => {
        const score = {}, classes = Object.keys(model);
        classes.forEach((c) => { score[c] = 0; });
        $('#input').val().toLowerCase().replaceAll('ё', 'е').split(/\s+/).forEach((wrd) => {
            wrd = `^${wrd}`; // BOW marker
            for (let i = 0; i < wrd.length - 2; i++) {
                let t = wrd.substr(i, 3);
                classes.forEach((c) => { score[c] += model[c][t] || 0; });
            }
        });
        const chart = classes.map((c) => [c, score[c]]);
        chart.sort((a, b) => b[1] - a[1]);
        $('#result').text(chart.map(([t, s]) => `${t} (${s})`).join('\n'));
    };
    const ax = {'method': 'GET', 'dataType': 'json'};
    $.ajax({...ax, 'url': 'model.json'}).done((m) => { Object.assign(model, m); update(); });
    $.ajax({...ax, 'url': 'data.json'}).done((m) => { $('#words').text(Object.keys(m).map((c) => c + m[c].map((x) => `\n   ${x}`).join('')).join('\n')); });
    $('#input').on('input', update);
    $('#input').focus();
});