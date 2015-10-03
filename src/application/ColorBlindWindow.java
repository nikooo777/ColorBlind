package application;

import java.util.List;

import javafx.geometry.Insets;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.layout.BorderPane;
import javafx.scene.layout.VBox;
import javafx.scene.paint.Color;
import javafx.scene.shape.Circle;
import javafx.stage.Modality;
import javafx.stage.Screen;
import javafx.stage.Stage;

public class ColorBlindWindow
{
	// objects of the window
	private final Stage window;
	private final Scene scene;
	private final BorderPane rootLayout, centerLayout;
	private final VBox verticalLayout;
	private final Button button, button2, button3;
	private Color prefColor;

	public ColorBlindWindow()
	{
		this.window = new Stage();

		// layouts
		this.rootLayout = new BorderPane();
		this.centerLayout = new BorderPane();
		this.verticalLayout = new VBox(20);

		// layout settings
		setupLayouts();

		// for the sake of fun, get the primary screen size and set the window
		// to be that size without "maximizing" it
		final Rectangle2D screen = Screen.getPrimary().getVisualBounds();

		// scene containing the root layout
		this.scene = new Scene(this.rootLayout, screen.getWidth(), screen.getHeight());

		// setup the window
		setupWindow(screen);

		// nodes
		this.button = new Button("Real Reverse CB");
		this.button2 = new Button("Confusion Reverse CB");
		this.button3 = new Button("Fix color for method 1");

		// event listeners
		initListeners();

		this.verticalLayout.getChildren().addAll(this.button, this.button2, this.button3);

	}

	private void initListeners()
	{
		this.button.setOnAction(e ->
		{
			// clear the screen
			this.centerLayout.getChildren().clear();

			// generate all circles
			final List<Circle> circles = CircleFactory.generateCircles(this.centerLayout);

			// debug message
			System.out.println("Successfully generated all circles!");

			// change the color to the circles that have to pop
			if (this.prefColor != null)
			{
				CircleFactory.drawSecrets(CircleFactory.shapeCircle(this.scene), circles, this.prefColor);
				System.out.println("Using special color");
			} else
				CircleFactory.drawSecrets(CircleFactory.shapeCircle(this.scene), circles, null);

			// add all the circles to the layout
			this.centerLayout.getChildren().addAll(circles);
		});
		this.button2.setOnAction(e -> System.out.println("Sorry this is not yet implemented :("));

		this.button3.setOnAction(e ->
		{
			this.prefColor = ColorFactory.colorPicker();
		});
	}

	private void setupLayouts()
	{
		this.rootLayout.setRight(this.verticalLayout);
		this.rootLayout.setCenter(this.centerLayout);
		this.verticalLayout.setPadding(new Insets(20));
	}

	private void setupWindow(final Rectangle2D screen)
	{
		// set Stage boundaries to visible bounds of the main screen
		this.window.setTitle("Magic display");
		this.window.initModality(Modality.APPLICATION_MODAL);
		this.window.setX(screen.getMinX());
		this.window.setY(screen.getMinY());
		this.window.setWidth(screen.getWidth());
		this.window.setHeight(screen.getHeight());
		this.window.setScene(this.scene);
	}

	public void show()
	{
		// pop up the window
		this.window.showAndWait();
	}

}
