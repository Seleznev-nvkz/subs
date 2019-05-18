function Word(props) {
    return (
        <span onClick={(e) => {if (e.ctrlKey || e.metaKey || props.isMobile) props.onClick(props.value)}}>{props.value}</span>
    )
}

class Container extends React.Component {
    constructor(props) {
        props.isMobile = window.innerWidth <= 500;
        super(props);
    }

    renderWord(word) {
        const value = word.translation ? `${word.word}\n${word.translation}` : word.word;
        return (
            <Word value={value} onClick={(e)=>this.props.removeWord(e, word)}/>
        )
    }

    render() {
        return (
            <div className="container">
                {this.props.entries.map(value => this.renderWord(value))}
            </div>
        )
    }
}

class Subtitles extends React.Component {
    renderSubtitle(subtitle) {
        return (
            <li><span onClick={() => this.props.loadSubtitle(subtitle)}>{subtitle.name}</span>
                <span className="remove" onClick={() => this.props.removeSubtitle(subtitle)}>&nbsp;X</span></li>
        )
    }

    render() {
        return (
            <ul className="subtitles-container">
                {this.props.subtitles.map(value => this.renderSubtitle(value))}
            </ul>
        )
    }
}

class Help extends React.Component {
    render() {
        return (
            <div className="help pure-form">
                <fieldset>
                    <ul>
                        <li>1. File{this.props.filename? ` (${this.props.filename})`: ''}</li>
                        <li>2. Click on word to remove as known</li>
                        <li>3. Translate</li>
                    </ul>
                </fieldset>
            </div>
        )
    }
}

class Main extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            ready: true,
            subtitles: [],
            selected: {name: "", words: []},
        };

        this.translate = this.translate.bind(this);
        this.removeWord = this.removeWord.bind(this);
        this.loadSubtitle = this.loadSubtitle.bind(this);
        this.removeSubtitle = this.removeSubtitle.bind(this);
    }

    componentDidMount() {
        this.loadSubtitles()
    }

    loadSubtitles() {
        fetch('api/subtitles/', {method: 'GET'})
            .then(res => res.json())
            .then(res => this.setState({...this.state, subtitles: res}))
            .catch(error => console.error('Error:', error));
    }

    removeWord(e, entry) {
        fetch(`api/words/${entry.word}`, {method: 'DELETE'})
            .then(res => (res.ok) ? res : new Error(res.text()))
            .then(this.setState({...this.state, selected: {...this.state.selected,
                    words: this.state.selected.words.filter(i => i.word !== entry.word)}}))
            .catch(error => console.error('Error:', error));
    }

    loadSubtitle(subtitle) {
        this.setState({...this.state, selected: subtitle});
    }

    removeSubtitle(subtitle) {
        fetch(`api/subtitles/${subtitle.id}`, {method: 'DELETE'})
            .then(res => (res.ok) ? res : new Error(""))
            .then(this.setState({...this.state, subtitles: this.state.subtitles.filter(i => i.id !== subtitle.id)}))
            .catch(error => console.error('Error:', error));
    }

    upload(file) {
        let form = new FormData();
        form.append("file", file);
        fetch('api/words/upload', {method: 'POST', body: form})
            .then(res => res.json())
            .then(res => this.setState({...this.state, selected: res,
                subtitles: [...this.state.subtitles, res]}))
            .catch(error => console.error('Error:', error));
    }

    translate() {
        this.setState({...this.state, ready: false});
        fetch(`api/subtitles/${this.state.selected.id}/translate`, {method: 'GET'})
            .then(res => {this.setState({...this.state, ready: true}); return res.json()})
            .then(res => this.setState({...this.state, selected: res}))
            .catch(error => console.error('Error:', error));
    }

    getInputRef = (node) => {this.fileInput = node};
    setLocalFile = () => {this.fileInput.click()};

    render() {
        const buttonCls = "pure-button pure-button-primary";

        return (
            <div className={this.state.ready ? "" : "disabled"}>
                <div className="pure-form">
                    <fieldset>
                        <button className={buttonCls} onClick={this.setLocalFile}>Set local file</button>
                        <input type="file" onChange={event => {
                            this.upload(event.target.files[0])
                        }} ref={this.getInputRef} className="hidden"/>

                        <button className={buttonCls} onClick={this.translate} disabled={!this.state.ready}>Translate</button>
                    </fieldset>
                </div>
                <Help filename={this.state.selected.name}/>
                <Subtitles subtitles={this.state.subtitles} loadSubtitle={this.loadSubtitle} removeSubtitle={this.removeSubtitle}/>
                <Container entries={this.state.selected.words} removeWord={this.removeWord}/>
            </div>
        );
    }
}

ReactDOM.render(<Main/>, document.getElementById("root"));
