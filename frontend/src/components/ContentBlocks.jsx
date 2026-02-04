export default function ContentBlocks({ content }) {
  const lines = (content || "").split("\n");
  const blocks = [];
  let inCode = false;
  let codeBuffer = [];

  const flushCode = () => {
    if (codeBuffer.length) {
      blocks.push({ type: "code", text: codeBuffer.join("\n") });
      codeBuffer = [];
    }
  };

  for (const line of lines) {
    if (line.trim().startsWith("```") && !inCode) {
      inCode = true;
      continue;
    }
    if (line.trim().startsWith("```") && inCode) {
      inCode = false;
      flushCode();
      continue;
    }

    if (inCode) {
      codeBuffer.push(line);
      continue;
    }

    if (!line.trim()) {
      blocks.push({ type: "spacer" });
      continue;
    }

    if (line.startsWith("### ")) {
      blocks.push({ type: "h3", text: line.replace("### ", "") });
      continue;
    }

    if (line.startsWith("- ")) {
      blocks.push({ type: "li", text: line.replace("- ", "") });
      continue;
    }

    blocks.push({ type: "p", text: line });
  }

  if (codeBuffer.length) flushCode();

  return (
    <div className="text-sm text-slate-700">
      {blocks.map((b, i) => {
        if (b.type === "spacer") return <div key={i} className="h-2" />;
        if (b.type === "h3") return <h3 key={i} className="mt-4 text-base font-semibold text-slate-900">{b.text}</h3>;
        if (b.type === "li") return <div key={i} className="ml-4 list-item list-disc text-slate-600">{b.text}</div>;
        if (b.type === "code") {
          return (
            <pre key={i} className="mt-3 overflow-x-auto rounded-xl bg-slate-900 p-4 text-xs text-slate-100">
              <code>{b.text}</code>
            </pre>
          );
        }
        return <p key={i} className="mb-2">{b.text}</p>;
      })}
    </div>
  );
}
