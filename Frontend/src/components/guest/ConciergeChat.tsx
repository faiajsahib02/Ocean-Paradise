import React, { useState, useRef, useEffect } from 'react';
import axios from 'axios';

type Message = {
  role: 'user' | 'bot';
  text: string;
};

export const ConciergeChat = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<Message[]>([
    { role: 'bot', text: "Hello! I'm the AI Concierge. Ask me anything about the hotel!" }
  ]);
  const [loading, setLoading] = useState(false);
  
  // Auto-scroll to bottom
  const endRef = useRef<HTMLDivElement>(null);
  useEffect(() => endRef.current?.scrollIntoView({ behavior: "smooth" }), [messages]);

  const sendMessage = async () => {
    if (!input.trim()) return;

    // 1. Add User Message
    const userMsg = input;
    setMessages(prev => [...prev, { role: 'user', text: userMsg }]);
    setInput("");
    setLoading(true);

    try {
      // 2. Call Backend API
      const res = await axios.post('http://localhost:8081/api/rag/ask', {
        question: userMsg
      });

      // 3. Add Bot Response
      setMessages(prev => [...prev, { role: 'bot', text: res.data.answer }]);
    } catch (err) {
      setMessages(prev => [...prev, { role: 'bot', text: "Sorry, I'm having trouble connecting to the front desk." }]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end font-sans">
      
      {/* The Chat Window */}
      {isOpen && (
        <div className="bg-white w-80 h-96 rounded-lg shadow-xl mb-4 flex flex-col border border-gray-200 overflow-hidden">
          {/* Header */}
          <div className="bg-blue-600 text-white p-3 font-bold flex justify-between items-center">
            <span>🤖 Smart Concierge</span>
            <button onClick={() => setIsOpen(false)} className="text-white hover:text-gray-200">✕</button>
          </div>

          {/* Messages Area */}
          <div className="flex-1 p-4 overflow-y-auto bg-gray-50">
            {messages.map((m, i) => (
              <div key={i} className={`mb-3 flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                <div className={`max-w-[80%] p-2 rounded-lg text-sm ${
                  m.role === 'user' 
                    ? 'bg-blue-500 text-white rounded-br-none' 
                    : 'bg-white border border-gray-200 text-gray-800 rounded-bl-none shadow-sm'
                }`}>
                  {m.text}
                </div>
              </div>
            ))}
            {loading && <div className="text-gray-400 text-xs ml-2">Thinking...</div>}
            <div ref={endRef} />
          </div>

          {/* Input Area */}
          <div className="p-3 bg-white border-t border-gray-100 flex gap-2">
            <input
              className="flex-1 border border-gray-300 rounded px-3 py-1 text-sm focus:outline-none focus:border-blue-500"
              placeholder="Ask about pool, parking..."
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
            />
            <button 
              onClick={sendMessage}
              disabled={loading}
              className="bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700 disabled:bg-gray-300"
            >
              ➤
            </button>
          </div>
        </div>
      )}

      {/* The Floating Button (FAB) */}
      {!isOpen && (
        <button
          onClick={() => setIsOpen(true)}
          className="bg-blue-600 text-white w-14 h-14 rounded-full shadow-lg flex items-center justify-center text-2xl hover:bg-blue-700 transition transform hover:scale-105"
        >
          💬
        </button>
      )}
    </div>
  );
};
