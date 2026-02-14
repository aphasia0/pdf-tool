#!/usr/bin/env python3
"""
Simple script to create a test encrypted PDF using reportlab
"""
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import letter
import PyPDF2
import sys

def create_test_pdf():
    # Create a simple PDF
    temp_pdf = "temp_unencrypted.pdf"
    encrypted_pdf = "test_encrypted.pdf"
    password = "testpassword123"
    
    # Create PDF with reportlab
    c = canvas.Canvas(temp_pdf, pagesize=letter)
    c.setFont("Helvetica", 12)
    
    # Add some content
    c.drawString(100, 750, "Test PDF Document")
    c.drawString(100, 730, "")
    c.drawString(100, 710, "This is a test PDF that will be encrypted with a password.")
    c.drawString(100, 690, "")
    c.drawString(100, 670, "Test content:")
    c.drawString(120, 650, "- Line 1: Sample data")
    c.drawString(120, 630, "- Line 2: More sample data")
    c.drawString(120, 610, "- Line 3: Additional content")
    c.drawString(100, 590, "")
    c.drawString(100, 570, f"Password for this PDF: {password}")
    
    c.save()
    
    # Encrypt the PDF
    pdf_reader = PyPDF2.PdfReader(temp_pdf)
    pdf_writer = PyPDF2.PdfWriter()
    
    # Add all pages
    for page in pdf_reader.pages:
        pdf_writer.add_page(page)
    
    # Encrypt with password
    pdf_writer.encrypt(password)
    
    # Write encrypted PDF
    with open(encrypted_pdf, 'wb') as output_file:
        pdf_writer.write(output_file)
    
    # Clean up temp file
    import os
    os.remove(temp_pdf)
    
    print(f"Created encrypted test PDF: {encrypted_pdf}")
    print(f"Password: {password}")

if __name__ == "__main__":
    try:
        create_test_pdf()
    except ImportError as e:
        print(f"Error: Missing required library. Install with:")
        print("pip3 install reportlab PyPDF2")
        sys.exit(1)
