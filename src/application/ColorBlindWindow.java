package application;

import java.awt.Button;

import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.layout.BorderPane;
import javafx.scene.layout.VBox;
import javafx.stage.Screen;
import javafx.stage.Stage;

public class ColorBlindWindow
{
	private final Stage window;
	private Scene scene;

	public ColorBlindWindow()
	{
		this.window = new Stage();
		this.window.setTitle("Magic display");

		// layouts
		final BorderPane rootLayout = new BorderPane();
		final VBox verticalLayout = new VBox(20);

		// layout settings
		rootLayout.setRight(verticalLayout);

		// for the sake of fun, get the primary screen size and set the window to be that size without "maximizing" it
		final Rectangle2D screen = Screen.getPrimary().getVisualBounds();

		// scene containing the root layout
		final Scene scene = new Scene(rootLayout, screen.getWidth(), screen.getHeight());

		// set Stage boundaries to visible bounds of the main screen
		this.window.setX(screen.getMinX());
		this.window.setY(screen.getMinY());
		this.window.setWidth(screen.getWidth());
		this.window.setHeight(screen.getHeight());

		this.window.setScene(scene);

		// nodes
		final Button button = new Button("Real Reverse CB");
		final Button button2 = new Button("Confusion Reverse CB");
	}

}
