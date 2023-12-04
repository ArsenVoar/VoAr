VoAr - Voice of Articles Overview: VoAr is a web application written in Go that serves as a platform for creating, managing, and sharing articles. It integrates Google authentication for user registration and includes features such as article creation, listing, and detailed article views. The project structure follows best practices in Go programming and incorporates improvements, error handling, and comments for better code readability and maintainability.

Previous Repositories: This project consolidates previous versions, bringing together improvements and features from earlier repositories. The transition aims to streamline development and provide a unified codebase for future enhancements.

Features:

User Authentication: Utilizes Google OAuth2 for user registration and authentication. Article Management: Allows users to create, view, and list articles. Improved Code Quality: Includes enhanced error handling, comments, and addressing weak points in the codebase. Structured Project: Follows a well-organized project structure for clarity and scalability.

How to Run the Project:

Clone the Repository:

bash | Copy code | -git clone https://github.com/your-username/VoAr.git | cd VoAr |

Set Up Environment Variables:

Create a .env file (and name it st.env) in the project root.sql Add the necessary environment variables for Google OAuth credentials, etc. For PostgreSQL database you can do your own or import from file mydb. Install Dependencies:

bash | Copy code | go get -u ./... |

Build and Run:

bash | Copy code | go build -o VoAr | ./VoAr |

Access the Application:

Open your web browser and navigate to http://localhost:8080.

Contributing:

Feel free to contribute by submitting issues, feature requests, or pull requests. Your input is valuable for the continuous improvement of VoAr.
