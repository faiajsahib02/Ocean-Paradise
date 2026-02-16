from reportlab.lib.pagesizes import letter
from reportlab.pdfgen import canvas
import os

# Create uploads directory if it doesn't exist
os.makedirs("./Backend/uploads", exist_ok=True)

# Create PDF
pdf_path = "./Backend/uploads/policy.pdf"
c = canvas.Canvas(pdf_path, pagesize=letter)
c.setFont("Helvetica", 12)
c.drawString(100, 750, "Hotel Policy Document")
c.drawString(100, 720, "")
c.drawString(100, 700, "The secret password for free coffee is 'Flamingo'.")
c.drawString(100, 680, "")
c.drawString(100, 660, "All guests who mention this password to the concierge")
c.drawString(100, 640, "will receive a complimentary cappuccino.")
c.save()

print(f"PDF created successfully at: {pdf_path}")
