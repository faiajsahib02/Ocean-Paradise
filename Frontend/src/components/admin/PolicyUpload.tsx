import React, { useState } from 'react';
import axios from 'axios';

export const PolicyUpload = () => {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<string>("");
  const [loading, setLoading] = useState(false);

  const handleUpload = async () => {
    if (!file) return;

    setLoading(true);
    setStatus("Reading PDF and training AI... (This may take a moment)");

    const formData = new FormData();
    formData.append('pdf', file); // Must match the backend key "pdf"

    try {
      await axios.post('http://localhost:8081/api/rag/ingest', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      setStatus("✅ Success! The Concierge has learned the new rules.");
    } catch (err) {
      setStatus("❌ Error: Failed to train the AI.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 bg-white rounded-lg shadow-md max-w-md">
      <h2 className="text-xl font-bold mb-4">🧠 Update Concierge Knowledge</h2>
      <p className="text-gray-600 mb-4 text-sm">
        Upload a new Hotel Policy PDF. The AI will read it instantly and answer guest questions based on it.
      </p>

      <input 
        type="file" 
        accept=".pdf"
        onChange={(e) => setFile(e.target.files?.[0] || null)}
        className="block w-full text-sm text-slate-500 mb-4 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
      />

      <button
        onClick={handleUpload}
        disabled={!file || loading}
        className={`w-full py-2 px-4 rounded font-bold text-white transition ${
          loading ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'
        }`}
      >
        {loading ? "Training AI..." : "Upload & Train"}
      </button>

      {status && (
        <div className={`mt-4 p-3 rounded text-sm ${status.includes("Error") ? "bg-red-50 text-red-700" : "bg-green-50 text-green-700"}`}>
          {status}
        </div>
      )}
    </div>
  );
};
