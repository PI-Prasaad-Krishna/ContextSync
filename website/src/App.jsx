import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';

function Navigation() {
  const location = useLocation();
  
  return (
    <nav className="navbar animate-up">
      <Link to="/" className="logo">ContextSync</Link>
      <div className="nav-links">
        {location.pathname === '/' ? (
          <Link to="/docs">Docs</Link>
        ) : (
          <Link to="/">Home</Link>
        )}
        <a href="https://github.com/PI-Prasaad-Krishna/ContextSync">GitHub</a>
      </div>
    </nav>
  );
}

function Home({ downloads, releaseData }) {
  return (
    <>
      <section className="hero animate-up">
        <div className="hero-content">
          <h1>Perfect Memory For Your AI.</h1>
          <p>
            ContextSync is a zero-dependency background daemon that maintains a hyper-dense, 
            continuously updated memory bank for AI coding agents. Stop wasting tokens on 
            full-repo scans and eliminate the "Lost in the Middle" effect.
          </p>
          <div className="hero-buttons">
            <a href="#install" className="brutalist-btn">Download Now</a>
            <Link to="/docs" className="brutalist-btn secondary">Read Docs</Link>
          </div>
        </div>

        <div className="mock-terminal">
          <div className="terminal-header">
            <div className="terminal-dot dot-red"></div>
            <div className="terminal-dot dot-yellow"></div>
            <div className="terminal-dot dot-green"></div>
          </div>
          <div className="terminal-body">
            <span className="comment"># 1. Initialize in your project</span><br/>
            <span className="command">~</span> ./contextsync init<br/><br/>
            
            <span className="comment"># 2. Start the daemon in the background</span><br/>
            <span className="command">~</span> ./contextsync watch<br/>
            <span className="comment">ContextSync daemon is now watching...</span><br/><br/>

            <span className="comment"># 3. Edit code. ContextSync handles the rest.</span><br/>
            <span className="command">~</span> <span style={{color: '#00e676'}}>Successfully synced files out=.context.md</span>
          </div>
        </div>
      </section>

      <section id="install" className="animate-up delay-1" style={{marginTop: '4rem'}}>
        <h2>Installation Packages</h2>
        <p>Select your operating system.</p>
        
        <div className="grid-3" style={{marginTop: '2rem'}}>
          <div className="brutalist-card">
            <h3>Windows</h3>
            <p style={{marginBottom: '1.5rem'}}>64-bit Executable (.exe)</p>
            <a href={downloads.windows} className="brutalist-btn" style={{width: '100%', textAlign: 'center'}}>Download</a>
          </div>
          
          <div className="brutalist-card">
            <h3>macOS</h3>
            <p style={{marginBottom: '1.5rem'}}>Apple Silicon & Intel</p>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
              <a href={downloads.macArm} className="brutalist-btn" style={{textAlign: 'center', backgroundColor: '#000', color: '#fff'}}>Apple Silicon (M1/M2)</a>
              <a href={downloads.macIntel} className="brutalist-btn secondary" style={{textAlign: 'center'}}>Intel Mac</a>
            </div>
          </div>

          <div className="brutalist-card">
            <h3>Linux</h3>
            <p style={{marginBottom: '1.5rem'}}>64-bit Binary</p>
            <a href={downloads.linux} className="brutalist-btn" style={{width: '100%', textAlign: 'center', backgroundColor: '#000', color: '#fff'}}>Download</a>
          </div>
        </div>
      </section>
    </>
  );
}

function Docs() {
  return (
    <div className="docs-page animate-up">
      <h2>Quick Start</h2>
      <p>ContextSync is a standalone binary. Once downloaded:</p>
      <ol style={{marginBottom: '1.5rem', marginLeft: '2rem'}}>
        <li>Extract the downloaded archive (<code>.zip</code> or <code>.tar.gz</code>).</li>
        <li>Place the extracted <code>contextsync</code> file in your project folder (or system PATH).</li>
      </ol>

      <p>Initialize a boilerplate context file in your project root. This creates a highly optimized markdown file that you can feed directly to your AI agent.</p>
      <pre><code>./contextsync init</code></pre>

      <p>Start the file-watching daemon. It will listen for file saves, debounce them, and instantly inject the exact diffs into your `.context.md` file.</p>
      <pre><code>./contextsync watch</code></pre>

      <h2>JSON Configuration</h2>
      <p>You can override the default daemon behavior by creating a <code>.contextsync.json</code> file in your project root.</p>
      <pre><code>{`{
  "debounce": "2s",
  "max_events": 15,
  "out": ".context.md",
  "debug": false
}`}</code></pre>
      <ul>
        <li><strong>debounce:</strong> How long to wait after a save before syncing (prevents spamming on Ctrl+S).</li>
        <li><strong>max_events:</strong> The maximum number of historical file changes to keep in memory before pruning.</li>
      </ul>

      <h2>Hybrid Memory Architecture</h2>
      <p>Because ContextSync only records your most <em>recent</em> file diffs, the AI will know exactly what you just did, but it might lose sight of the big picture.</p>
      
      <div className="brutalist-card" style={{marginTop: '2rem'}}>
        <h3>The Base Context Pipeline</h3>
        <p>To solve this, ContextSync features a native <strong>Base Context</strong> pipeline:</p>
        <ul>
          <li><strong>Long-Term Memory:</strong> Write a static <code>Architecture.md</code> file that explains your overall project goals, structs, and patterns.</li>
          <li><strong>Configuration:</strong> Add <code>"base_context_file": "Architecture.md"</code> to your JSON config.</li>
          <li><strong>The Magic:</strong> On every file save, ContextSync dynamically copies your entire architecture and glues it to the very top of <code>.context.md</code>, right above the rolling diffs.</li>
        </ul>
      </div>
    </div>
  );
}

function Privacy() {
  return (
    <div className="docs-page animate-up">
      <h2>Privacy Policy</h2>
      <p style={{marginBottom: '3rem'}}><strong>Effective Date:</strong> {new Date().toLocaleDateString()}</p>
      
      <h3>1. Fully Local Operation</h3>
      <p>ContextSync is designed with absolute privacy in mind. It operates 100% locally on your machine. The daemon monitors your local file system and writes directly to your local <code>.context.md</code> file.</p>
      
      <h3>2. Zero Data Collection</h3>
      <p>We do not collect, store, transmit, or monitor any of your data. There are no analytics, no telemetry, and no crash reporting built into ContextSync.</p>
      
      <h3>3. No Cloud Communication</h3>
      <p>ContextSync makes zero outbound network requests regarding your code. The only network request made is by this documentation website to fetch the latest GitHub releases for the download buttons.</p>

      <h3>4. Your Code is Yours</h3>
      <p>Your source code, file structures, and context files never leave your computer. You have complete ownership and control over the data generated by ContextSync.</p>
    </div>
  );
}

function App() {
  const [releaseData, setReleaseData] = useState(null);
  const [downloads, setDownloads] = useState({
    windows: '#',
    macIntel: '#',
    macArm: '#',
    linux: '#'
  });

  useEffect(() => {
    fetch('https://api.github.com/repos/PI-Prasaad-Krishna/ContextSync/releases/latest')
      .then(res => res.json())
      .then(data => {
        setReleaseData(data);
        if (data.assets) {
          const links = { ...downloads };
          data.assets.forEach(asset => {
            if (asset.name.includes('windows') && asset.name.endsWith('.zip')) links.windows = asset.browser_download_url;
            if (asset.name.includes('darwin-amd64') && asset.name.endsWith('.tar.gz')) links.macIntel = asset.browser_download_url;
            if (asset.name.includes('darwin-arm64') && asset.name.endsWith('.tar.gz')) links.macArm = asset.browser_download_url;
            if (asset.name.includes('linux') && asset.name.endsWith('.tar.gz')) links.linux = asset.browser_download_url;
          });
          setDownloads(links);
        }
      })
      .catch(err => console.error("Error fetching release:", err));
  }, []);

  return (
    <BrowserRouter>
      <div className="app-container">
        <Navigation />

        <Routes>
          <Route path="/" element={<Home downloads={downloads} releaseData={releaseData} />} />
          <Route path="/docs" element={<Docs />} />
          <Route path="/privacy" element={<Privacy />} />
        </Routes>

        <footer className="animate-up delay-2">
          <p>&copy; {new Date().getFullYear()} ContextSync. All rights reserved.</p>
          <p>
            <a href="https://github.com/PI-Prasaad-Krishna/ContextSync">GitHub Repository</a> | <Link to="/privacy">Privacy Policy</Link>
          </p>
        </footer>
      </div>
    </BrowserRouter>
  )
}

export default App;
