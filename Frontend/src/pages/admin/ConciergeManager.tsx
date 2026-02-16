import { PolicyUpload } from '../../components/admin/PolicyUpload';

const ConciergeManager = () => {
  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">AI Concierge Manager</h1>
        <p className="text-gray-600">
          Manage the AI-powered concierge system. Upload hotel policies and the AI will instantly learn them 
          to answer guest questions 24/7.
        </p>
      </div>

      <div className="grid gap-6">
        {/* Policy Upload Section */}
        <section>
          <PolicyUpload />
        </section>

        {/* Info Section */}
        <section className="bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-blue-900 mb-3">How It Works</h3>
          <ul className="space-y-2 text-blue-800 text-sm">
            <li className="flex items-start">
              <span className="mr-2">📄</span>
              <span>Upload a PDF containing hotel policies, rules, or FAQs</span>
            </li>
            <li className="flex items-start">
              <span className="mr-2">🧠</span>
              <span>Our AI reads and understands the document using advanced language models</span>
            </li>
            <li className="flex items-start">
              <span className="mr-2">💬</span>
              <span>Guests can chat with the AI concierge on your website to get instant answers</span>
            </li>
            <li className="flex items-start">
              <span className="mr-2">🔄</span>
              <span>Update the knowledge base anytime by uploading a new policy document</span>
            </li>
          </ul>
        </section>

        {/* Stats Section */}
        <section className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Quick Tips</h3>
          <div className="grid md:grid-cols-2 gap-4 text-sm text-gray-600">
            <div className="flex items-start space-x-2">
              <span className="text-green-500 font-bold">✓</span>
              <span>Use clear, well-formatted PDFs for best results</span>
            </div>
            <div className="flex items-start space-x-2">
              <span className="text-green-500 font-bold">✓</span>
              <span>Include check-in/out times, amenities, and policies</span>
            </div>
            <div className="flex items-start space-x-2">
              <span className="text-green-500 font-bold">✓</span>
              <span>Update regularly when policies change</span>
            </div>
            <div className="flex items-start space-x-2">
              <span className="text-green-500 font-bold">✓</span>
              <span>Test the chat widget as a guest to verify accuracy</span>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
};

export default ConciergeManager;
