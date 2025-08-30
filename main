import tkinter as tk
from tkinter import filedialog, messagebox
import base64
import os

# Variables to store file and folder paths
file_path = None
output_folder = None

# Function to select a base64 encoded input file
def select_file():
    global file_path
    file_path = filedialog.askopenfilename(filetypes=[("Base64 Files", "*.ote")])
    if file_path:
        file_label.config(text=f"Selected File: {os.path.basename(file_path)}")

# Function to select the output folder
def select_output_folder():
    global output_folder
    output_folder = filedialog.askdirectory()
    if output_folder:
        folder_label.config(text=f"Output Folder: {output_folder}")

# Function to convert from base64 back to the original file content (UTF-8)
def convert_from_base64():
    if not file_path or not output_folder:
        messagebox.showerror("Error", "Please select both a file and an output folder.")
        return
    
    try:
        # Read the base64 encoded file content
        with open(file_path, 'r') as file:
            base64_string = file.read()

        # Decode the base64 string back to binary data
        file_content = base64.b64decode(base64_string)
 ### REWRITE!!!
        # Prepare the new file name with its original extension (assuming it was originally a .txt file)
        file_name = os.path.basename(file_path)
        original_file_name = os.path.splitext(file_name)[0]  # Removing the .ote extension
        new_file_path = os.path.join(output_folder, original_file_name)

        # Write the decoded content to a new file
        with open(new_file_path, 'wb') as new_file:
            new_file.write(file_content)

        # Display success message
        messagebox.showinfo("Success", f"File converted back and saved as:\n{new_file_path}")
    except Exception as e:
        # Handle errors
        messagebox.showerror("Error", f"An error occurred: {str(e)}")

# Function to create the GUI
def create_gui():
    global file_label, folder_label

    # Create the main window
    root = tk.Tk()
    root.title("OpenTunnelV2 File Converter")
    root.geometry("600x400")  # Increased window size

    # Create a label for "OpenTunnelV2" at the top-center
    label_title = tk.Label(root, text="OpenTunnelV2", font=("Helvetica", 18), anchor="center")
    label_title.pack(pady=1)
    label_subtitle = tk.Label(root, text="Host", font=("Helvetica", 12), anchor="center")
    label_subtitle.pack(pady=1)

    # Create a button to select the input Base64 file
    select_button = tk.Button(root, text="Select Base64 File", command=select_file)
    select_button.pack(pady=10)

    # Label to show the selected file
    file_label = tk.Label(root, text="Selected File: None", anchor="w")
    file_label.pack(pady=5)

    # Create a button to select the output folder
    folder_button = tk.Button(root, text="Select Folder", command=select_output_folder)
    folder_button.pack(pady=5)
    
    # Create a label for output folder selection
    folder_label = tk.Label(root, text="Output Folder: None", anchor="w")
    folder_label.pack(pady=5)
    
    # Create a "Go" button to trigger the conversion back to original file
    go_button = tk.Button(root, text="Go", command=convert_from_base64)
    go_button.pack(pady=20)

    # Start the Tkinter event loop
    root.mainloop()

# Main execution
if __name__ == "__main__":
    create_gui()
