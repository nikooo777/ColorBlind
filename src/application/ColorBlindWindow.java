package application;

import java.util.List;

import javafx.geometry.Insets;
import javafx.geometry.Rectangle2D;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.control.Menu;
import javafx.scene.control.MenuBar;
import javafx.scene.control.MenuItem;
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
	private final MenuItem item1;
	private final Menu menu1;
	private final MenuBar mb;

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
		this.button2.setDisable(true);
		this.button3 = new Button("Fix color for method 1");

		// menubar
		this.mb = new MenuBar();

		// menus
		this.menu1 = new Menu("File");

		// items
		this.item1 = new MenuItem("close");

		// setup menubar
		this.menu1.getItems().add(this.item1);
		this.mb.getMenus().add(this.menu1);

		// event listeners
		initListeners();

		this.rootLayout.setTop(this.mb);
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
		this.item1.setOnAction(e -> this.window.close());
	}

	private void setupLayouts()
	{
		this.rootLayout.setRight(this.verticalLayout);
		this.rootLayout.setCenter(this.centerLayout);
		this.verticalLayout.setPadding(new Insets(20));
		this.centerLayout.setPadding(new Insets(10));
		this.centerLayout.setStyle("-fx-background-color: #DADADA");
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
